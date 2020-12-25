# AKS へのデプロイ

**注意:** Git Bash から作業を行うと Azure CLI のバグによりコマンドが失敗します。Ubuntu や Mac OS を使用することを推奨します。

AKS へデプロイするために以下のツールが必要なため事前にインストールしてください。

- [Azure CLI]
- [jq]
- kubectl

このドキュメントでは [AKS] や Kubernetes に関する説明を行いませんので、公式ドキュメントを参照してください。

## AKS クラスタを作成

Azure に [AKS] クラスタを作成してください。
クラスタの作成手順は省略します。

## Azure にカスタムロールを追加

このプロダクトを実行するには [AKS] のノードプール (`AgentPool`) の取得・更新を行う権限を必要とします。

[Azure CLI] と [jq] を使用して Azure に最小権限のカスタムロールを登録します。
これらのツールをインストールしてください。

[Azure CLI] からログインしてください。
またサブスクリプション ID を取得します。

```bash
$ az login
$ SUBSCRIPTION_ID=$(az account show | jq -r .id)
$ echo $SUBSCRIPTION_ID
```

カスタムロール JSON のひな形を使用してカスタムロールを登録してください。

```bash
$ ROLE_JSON=$(sed -e "s/subscriptionID/$SUBSCRIPTION_ID/" ./deployments/azure/AzureKubernetesAgentPoolOperatorRole.json)
$ echo $ROLE_JSON
$ az role definition create --role-definition "$ROLE_JSON"
$ echo $?
0
```

## Azure にアプリを登録

このプロダクトから [AKS] ノードプールを更新するためのカスタムアプリを登録します。

対象の [AKS] クラスタのリソース ID を取得します。

```bash
$ RESOURCE_ID=$(az aks show --resource-group "<Your resource group name>" --name "<Your resource name>" | jq -r .id)
$ echo $RESOURCE_ID
```

アプリを登録します。

**注意:** `az ad sp create-for-rbac` で作成した時のパスワードは後ほど `AZURE_CLIENT_SECRET` として使用します。

```bash
$ APP_NAME=aks-scheduled-poolscaler
$ az ad sp create-for-rbac --name "$APP_NAME"
{
  "appId": "ffffffff-ffff-ffff-ffff-ffffffffffff",
  "displayName": "xxxxxxxx",
  "name": "http://xxxxxxxx",
  "password": "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
  "tenant": "ffffffff-ffff-ffff-ffff-ffffffffffff"
}
$ APP_ID=$(az ad sp show --id "http://$APP_NAME" --query appId | jq -r .)
$ echo $APP_ID
```

作成したアプリにカスタムロールを割り当てます。

**注意:** Azure CLI の [バグ](https://github.com/MicrosoftDocs/azure-docs/issues/24857) により Git Bash からは以下のコマンドが成功しません。

```bash
$ az role assignment create --assignee "$APP_ID" --role "Azure Kubernetes Agent Pool Operator Role" --scope "$RESOURCE_ID"
```

これで [AKS] クラスタ外の準備は完了です。

## デプロイ先の名前空間を作成

kubectl で対象のクラスタを操作するために対象の [AKS] リソースを指定して、認証情報を取得してください。

```bash
$ az aks get-credentials --resource-group "Your resource group name" --name "Your resource name"
```

このプロダクトをデプロイする名前空間を作成してください。

```bash
$ NAMESPACE=aks-scheduled-poolscaler
$ kubectl create namespace "$NAMESPACE"
```

## AKS を操作するための認証情報のシークレットを作成

AKS を操作するための認証情報のシークレットを作成します。
この認証情報は Azure SDK fo Go の
[環境変数ベースの認証情報の取得機能](https://docs.microsoft.com/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication)
にて参照されます。

以下のコマンドを実行して認証情報のシークレットを作成します。

```bash
$ TENANT_ID=$(az account show | jq -r .tenantId)
$ SUBSCRIPTION_ID=$(az account show | jq -r .id)
$ echo $TENANT_ID $SUBSCRIPTION_ID $APP_ID
$ APP_ID=$(az ad sp show --id "http://$APP_NAME" --query appId | jq -r .)
$ kubectl create secret generic aks-scheduled-poolscaler-client-credentials -n "$NAMESPACE" \
  --from-literal=tenantid="$TENANT_ID" \
  --from-literal=subscriptionid="$SUBSCRIPTION_ID" \
  --from-literal=clientid="$APP_ID" \
  --from-literal=clientsecret="Your client secret (app password)"
```

## 設定ファイルリソースを作成

プロダクトが変更する対象のリソース名やルールを記載した設定ファイルリソースを作成します。
[deployments/kubernetes/aks-scheduled-poolscaler-config-sample.yml](../deployments/kubernetes/aks-scheduled-poolscaler-config-sample.yml)
をベースに設定ファイルリソースを作成してください。
作成したリソース定義ファイル名が "aks-scheduled-poolscaler-config.yml" であるものとして続行します。

以下のコマンドを実行して設定ファイルリソースを Kubernetes に登録してください。

```bash
$ kubectl apply -f aks-scheduled-poolscaler-config.yml -n "$NAMESPACE"
```

## CronJob リソースを作成

以下のコマンドを実行して CronJob リソースを Kubernetes に登録してください。

```bash
$ kubectl apply -f deployments/kubernetes/aks-scheduled-poolscaler.yml -n "$NAMESPACE"
```

以上で AKS へのデプロイは完了です。

[aks]: https://azure.microsoft.com/services/kubernetes-service/
[azure cli]: https://docs.microsoft.com/cli/azure/
[jq]: https://stedolan.github.io/jq/

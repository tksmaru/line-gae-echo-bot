# line-gae-echo-bot

## 概要

発言に対してオウム返しするLINE bot。GAE/GO上で動作する。


## 利用ツール

プロジェクトベースのビルドツールとして[gb](https://getgb.io/)、及びGAE/GO用のプラグインとして[gb-gae](https://github.com/PalmStoneGames/gb-gae)を使用する。

```
go get github.com/constabulary/gb/...
go get code.palmstonegames.com/gb-gae
```

以降の説明はこれらの利用を前提とする。


## 事前準備

botを動かすための事前準備として以下が必要。

* LINEの[BOT API Trial Account](https://business.line.me/services/products/4/introduction)を作成していること
* Google Cloud Platformにプロジェクトを作成していること
* [Google App Engine SDK](https://cloud.google.com/appengine/downloads)をセットアップしておくこと


## 実行準備

#### 1) 本プロジェクトをcloneする

```
git clone git@github.com:tksmaru/line-gae-bot-sample.git
```

#### 2) 依存ライブラリの解決

以下のコマンドを実行して、本プロジェクトが依存する各種ライブラリを取得する。

```
cd /path/to/line-gae-echo-bot/
gb vendor restore
```

#### 3) `app.yaml`の編集

`src/app.yaml`の以下の項目を編集する。

| 設定項目 | 説明 |
|:----|:----|
| CHANNEL_ID | Channel IDの値 |
| CHANNEL_SECRET | Channel Secretの値 |
| MID | MIDの値 |


## 実行

### GAE/GO上で動かす

#### 1) デプロイ

プロジェクトのルートディレクトリで以下のコマンドを実行してGAE/GOにデプロイする。

```
cd /path/to/line-gae-echo-bot/
gb gae deploy --application=<your app id> --version <version> src/app.yaml
```

#### 2) コールバックURLの登録

LINE developers connect上でcallback用のURLを登録する。

```
https://<your app id>.appspot.com:443/callback
```

**ポイント**

* SSLのポート番号を明示する
* 実際に設定が反映されるまでには若干のタイムラグがある(体感30分)

#### 3) botを友達に登録して話しかける

自分の発言に対してbotがオウム返しで反応したら成功。


### ローカル環境で動かす

自分のPCなどでとりあえず動かしたい場合は、プロジェクトのルートディレクトリで以下のコマンドを実行する。

```
cd /path/to/line-gae-echo-bot/
gb gae serve src/app.yaml
```

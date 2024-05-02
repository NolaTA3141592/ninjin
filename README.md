# ninjin
slackに送信した内容をdiscordに移してログを取りたい！

## ngrokによるデバッグ
- ngrokのインストール(https://zenn.dev/claustra01/articles/c4d22f187943cf)
- ターミナルから`ngrok http 3000`を実行
- 実行すると出てくるForwardingのURLをコピーしておく(https:// ~ .ngrok-free.app)
- slackのappを作成
- slack apiのサイトからappの編集を行う。
- verification tokenをコピー&コードに貼り付け
- 画面左側からEvent Subscriptionを編集
- Enable Eventsをonにする
- Request URL に先ほどコピーしたngrokのURLを貼り付け
- このURLの末尾に`/slack/events`を追加
- ninjinを起動(ローカルでも、Dockerコンテナでも良い)
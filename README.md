# Concurrent
Concurrentは分散マイクロブログ基盤です。

## Motivation
Concurrentは、「セルフホストでお一人様インスタンス建てたい！けどローカルが1人ぼっちなのは寂しい...」という問題を解決するために生まれました。
個々のサーバーが所有しているタイムライン(mastodonやmisskeyで言うところのローカル)を別のサーバーから閲覧ないしは書き込みができます。
また、自分が閲覧しているタイムラインに対して、どのサーバーの持ち物であってもリアルタイムなイベントを得ることができます。

これにより、どのサーバーにいても世界は一つであるように、壁のないコミュニケーションが可能です。

## How it works
Concurrentでは公開鍵を用いて、ユーザーが発行するメッセージに署名を行います。

これにより、そのツイートがその秘密鍵の持ち主によって行われたことが誰でも検証できるようになります。

ConcurrentではユーザーのIDはConcurrentアドレス(cc-address)(例えば、`CC3E31b2957F984e378EC5cE84AaC3871147BD7bBF`)を用いて識別されます。

## インスタンスの立ち上げ方
docs/README.mdを参照

## Contributing
コードのPRは必ずissueでその可否のコンセンサスをとってからにしてください(せっかく作ってくれたPRをcloseするのは心が痛いので)。


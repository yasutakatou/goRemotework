# goRemotework
**Go言語製のリモートワーク支援ツール。リモートワークに不慣れなジュニアエンジニアと、適切に管理したいマネージャーとの橋渡しをします**

# リモートワークにおけるジュニアエンジニアをオンボーディングするための課題
## 証拠
### ジュニアエンジニアの視点
- ①日報等に日々の作業をまとめてはいるものの、実際には何にどれくらい時間を消費したのか正確に把握し辛いため、何をやってたのかの可視化が難しい。また、日報を書く時間自体が手間でルーズになりがちになる
- ②作業に詰まっている時に、どの程度の調査量で報連相すればいいのかの判断が育ってないため、自己解決しようとハマり続けてしまう
- ③時間の管理能力が育ってないため、カレンダーツールに予定を入れたり工夫はするものの、会議予定とカブったりでうまく時間のアラートを受けられない

これらの事実によりオフィス勤務経験の無い(または少ない)ジュニアエンジニアが、時間管理に不慣れなままリモートワークをすることで、結果、仕事力向上が停滞し、逆にマネージャーにより多くの労力がかかることになる

### マネージャーの視点
続いてジュニアエンジニアの視点の対義としてマネージャーには以下の労力がかかることとなる

- ①がもたらすもの　日報を読んでもスキルプランニングが難しく。KPI的に考えられない。また、サボりや仕事をしていない時間が無いかの管理が必要な点について性善説に基づくしかない
- ②がもたらすもの　オフィス勤務のように背中越しで、詰まってる雰囲気を察して、適切なタイミングで声掛けができないため、想定より時間をかけて作業をしてしまう
- ③がもたらすもの　オフィスでは簡単な声掛けで作業漏れを防げたものをチェックする手間が増える。また、指示したタスクに漏れが無いかも確認する必要がある

## 主張
これは双方の抱える人材育成系のエンジニアリング技術の問題ではなく、**リモートワークがオフィス勤務をベースにしているため、ジュニアエンジニアが慣れることに時間がかかり、シニアエンジニアが適切な指導、育成の機会を与えるのが難しいという点が根本的な問題**として潜在する。<br>
そのため発生している以下のアンチパターンを参考に、デメリットを最小に軽減する仕組みをもったアプリ、SaaS、ツールがジュニアエンジニア向けに提供されていないと思われる

### アンチパターン
- リモート監視ツールを導入し、キーストロークの保存と、インカメから仕事状況を録画し、離籍等を記録。ツールの出力からジュニアエンジニアへの適切な支援を行う

この方法は性悪説に基づいていると言える。心理的安全性の低下、過剰監視、マイクロマネジメント等のオフィス勤務であっても非推奨となる管理行為を実施している

- 全日オフィス勤務　**※メリットも大きいため、主張に基づいたデメリットの側面にフォーカスしている**

生産性の向上目的で出社率を増やす目的に組み込まれているため育成の観点では有効と思われる。しかし、シニアエンジニアやマネージャーの業務がリモートワークに最適化された後の状況では通勤時間を作業時間に当てたいという状況も考えられる

## 保障
### 望ましい状態
- ①について　出来るだけ透過的に個人差の無い可視化と平準的な数値化が出来ている
- ②について　マネージャー、シニアエンジニアのタイミングで過剰な思考の迷路状態をキャッチアップ出来る
- ③について　マネージャー、シニアエンジニアと合意した作業タイミングでタスクに取り掛かれる

**なぜならジュニアエンジニアへの支援と育成は、対人行為をインタフェースとした一種のインプットとアウトプットであり、エンジニアリング経験を積むために重要なものであるからだ**

# 機能
- ①**バックグラウンドで使用したアプリケーション情報を収集し、それぞれに消費した時間を計測する**
  - キーストロークの計測監視未満のレイヤーであるアプリ情報レベルで収集することで過剰監視となる事を防ぐ
- ②**特定のアプリケーションの操作時間を集計し、一定時間を超える場合に定義されたRunbookを実行し、アラートを自動発砲する**
  - 煮詰まった状況を自動的にキャッチアップできるようにする。また、マネジメント側がアラート時間を設定することで、報連相への心理的安全を高める
- ③**指定された時間や日時で通知することで、どのタイミングで何をすれば良いのか気づくことができる**
  - 会議のカレンダーと切り離す事でタスクの忘却を防ぐ。また、これもマネジメント側と相違の設定により、双方のタスクへのコミット感が醸成される

# 具体的な動作
**業務PCが必ず起動、シャットダウンが行われる前提を利用し、①～③の動作を一定間隔でループさせることでツールとして落とし込む**

- ①バックグラウンドで使用したアプリケーション情報を収集し、それぞれに消費した時間を計測する
  - フォアグラウンドにあるウィンドウのアプリケーションタイトルを取得し、その時間を類型したものを定期的にファイルに書き出す
  - 正規表現により、同じアプリでも別のタイトルになるものをシュリンクして同じ項目として集計できる仕組みにより、単一のタスクとして扱うことが出来る　※例えば開いている画面によってタイトルが変わるブラウザ等
- ②特定のアプリケーションの操作時間を集計し、一定時間を超える場合に定義されたRunbookを実行し、アラートを自動発砲する
  - ①同様に特定のアプリケーションをフォアグラウンドで動かしている時間が一定時間を超えた場合、定義されたコマンドを実行する
  - これも①同様に正規表現でシュリンクして同じ項目として集計できる仕組みを設ける
- ③指定された時間や日時で通知することで、どのタイミングで何をすれば良いのか気づくことができる
  - 正規表現により、特定の日時、または時間で定義されたコマンドを実行する

## 使い方
各設定ファイルに定義されたルールで動作します。よって以下のようなチームでの運用が望ましいでしょう

- PMやシニアエンジニアは「タスク集計定義」と「スケジュール通知」の設定ファイルを準備し、チームでファイル共有します
- Windows起動時に自動起動する設定にします [設定例](https://note.com/bright_clover112/n/ncd35e325b202)
- 出力される作業時間のファイルを定期的にアップロードするようにします
- 作業時間のファイルをPMが評価し、適切なマネジメントになるようにチューニングします

# 動作画面
- ①バックグラウンドで使用したアプリケーション情報を収集し、それぞれに消費した時間を計測する

フォアグラウンドで操作していたウィンドウで正規表現でマッチしたものを集計します

![image](https://github.com/user-attachments/assets/99e76027-42f7-4da7-ad1f-90607ef72baa)

- ②特定のアプリケーションの操作時間を集計し、一定時間を超える場合に定義されたRunbookを実行し、アラートを自動発砲する

stackoverflowで検索し続けていたのでアラートが出た例

![image](https://github.com/user-attachments/assets/3737c0ae-4e53-4d22-b437-d96529958674)

- ③指定された時間や日時で通知することで、どのタイミングで何をすれば良いのか気づくことができる

毎時45分になったらポモドーロテクニックで休憩しましょう、な例

![image](https://github.com/user-attachments/assets/6fae4eee-fb2b-4451-a07f-9d65e9a2909e)

# インストール方法
ツールをこのリポジトリからもってきます

```
go get github.com/yasutakatou/goRemotework
```

それかクローンしてビルドするか

```
git clone https://github.com/yasutakatou/goRemotework
cd caplint
go build .
```

面倒なら[ここにバイナリファイル置いておくので](https://github.com/yasutakatou/goRemotework/releases)手元で展開するでもOKです

# アンインストール方法

Go言語なのでバイナリファイル消してあげればOK！

# 設定ファイル

**設定ファイルはtsv(tab split value)形式です。タブ区切りで定義します**

## タスク集計定義(と、アラートRunbook)

この設定ファイルでは集計するタスク名と、(設定する場合)時間超過のアラートを定義する

```
①定義名 ②集計するウィンドウ名の正規表現 ③時間超過(秒)　④超過時に実行されるコマンド  ⑥「⑤」のコマンドに与える引数
```

例えばAWSのマネージメントコンソールを眺めている時間を計測と、30分を越えるとアラートするのであれば以下のように設定します
(アラートの設定をPMやシニアエンジニアが行うので、エスカレーションへの障壁がさがります)

```
AWS	.*Amazon Web Services.*	1800	popup.bat	AWS_USe
```

アラートを出さないでウィンドウ名と類型時間だけ計測する場合は以下のように③～⑤に「NO」を書きます

```
メモ帳	.*メモ帳.*	NO	NO	NO
```

また、正規表現に当てはまらないウィンドウ名は全て「OTHER」にまとめられます。作業用BGMでYoutubeを開いた、といったものまで過剰監視しない作りになっています

## スケジュール通知

この設定ファイルでは設定した時間でアラートを実施します
タスク集計定義とファイルを分けているのは、チーム共有のアラート設定、プラス、各自のアラートを追記する事で個別のアラートにも対応できるために分けています

```
①アラートを出す正規表現 ④実行されるコマンド  ⑥「⑤」のコマンドに与える引数
```

日時は以下のフォーマットで記録されます

```
2025/01/06 13:45:38 Mon
```

年、月、日、時刻、曜日です。例えば毎日10～13時にアラートを出したい場合は以下のように設定します

```
.*/.*/.* 1[0-3]:.*:.* .*	popup.bat	10-13HourAlert
```

月曜日の9:00であれば以下です

```
.*/.*/.* 9:00:00 Mon	popup.bat	MondayAlert
```

後述のループ間隔のオプションを広げ過ぎると、アラート出せない可能性があるので注意してください

# オプション

```
  -debug
        [-debug=debug mode (true is enable)]
  -log
        [-log=logging mode (true is enable)]
  -loop int
        [-loop=incident check loop time (Seconds). ] (default 60)
  -outputconfig string
        [-outputconfig=specify the output file of the work history.] (default "output.txt")
  -scheduleconfig string
        [-scheduleconfig=specify the configuration file for scheduled alerts.] (default "schedule.ini")
  -tasksconfig string
        [-tasksconfig=specify the task aggregation config file.] (default "tasks.ini")

```

## -debug

デバッグモードで動作します。色々出力されます

## -log

デバッグモードで出たログを出力するオプションです

## -loop int

全体の動作をループさせる間隔(秒)です<br>
長すぎると正確に計測できないのでできるだけ短く、数秒単位で運用する方が良いでしょう。また、**スケジュール通知設定がループ間隔によって跨いでしまう場合は、アラートが出なくなってしまうのでその点の注意も必要です**

## -outputconfig string

ウィンドウ名と類型時間だけ計測するファイルをデフォルトの　output.txt　から変えたいときにこれでコンフィグファイル名を指定します<br>
OSシャットダウン時にアップロード、同期する前提のため、常に計測ファイルは上書きされます。よって、不意のOS再起動が考えられる場合はこのファイル名を日時等に変更した方が良いでしょう

## -scheduleconfig string

スケジュール通知設定ファイルをデフォルトの　schedule.ini　から変えたいときにこれでコンフィグファイル名を指定します

## -tasksconfig string

タスク集計定義設定ファイルをデフォルトの　tasks.ini　から変えたいときにこれでコンフィグファイル名を指定します

# ライセンス
BSD 3-Clause

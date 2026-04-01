## 2026-04-01 17:30 — v1 函數重命名 + v2 derivatives bug 修正 + v2 新增 31 個 endpoint

**改動摘要：** 將所有 v1 exported 函數加上 V1 前綴（37 個函數），修正 derivatives.go 的 permission 與 type assertion bug，新增 31 個 v2 REST API endpoint

**修改的檔案：**

### v2 bug 修正
- `v2/rest/derivatives.go` — `PermissionRead` → `PermissionWrite`；type assertion 改為 comma-ok + float64

### v1 函數重命名（37 個函數加上 V1 前綴）
- `v1/account.go` — Info→V1Info, KeyPermission→V1KeyPermission, Summary→V1Summary
- `v1/balances.go` — All→V1All
- `v1/credits.go` — All→V1All
- `v1/deposit.go` — New→V1New
- `v1/history.go` — Balance→V1Balance, Movements→V1Movements, Trades→V1Trades
- `v1/lendbook.go` — Get→V1Get, Lends→V1Lends
- `v1/margin_funding.go` — NewLend→V1NewLend, NewLoan→V1NewLoan, Cancel→V1Cancel, Status→V1Status, Credits→V1Credits, Offers→V1Offers
- `v1/margin_info.go` — All→V1All
- `v1/offers.go` — New→V1New, Cancel→V1Cancel, Status→V1Status
- `v1/order_book.go` — Get→V1Get
- `v1/orders.go` — All→V1All, CancelAll→V1CancelAll, Create→V1Create, Cancel→V1Cancel, CreateMulti→V1CreateMulti, CancelMulti→V1CancelMulti, Replace→V1Replace, Status→V1Status
- `v1/pairs.go` — All→V1All, AllDetailed→V1AllDetailed
- `v1/positions.go` — All→V1All, Claim→V1Claim
- `v1/stats.go` — All→V1All
- `v1/ticker.go` — Get→V1Get
- `v1/trades.go` — All→V1All
- `v1/wallet.go` — Transfer→V1Transfer, WithdrawCrypto→V1WithdrawCrypto, WithdrawWire→V1WithdrawWire
- 對應 16 個 test 檔案 + 3 個 integration test + 1 個 example 的 call site 同步更新

### v2 新增 endpoint（6 個新檔案 + 5 個既有檔案擴充）
- 新建 `v2/rest/account.go`（6 methods）、`movements.go`（2）、`alerts.go`（3）、`settings.go`（3）、`liquidations.go`（1）、`rankings.go`（1）
- 擴充 `positions.go`（+5）、`funding.go`（+4）、`derivatives.go`（+1）、`status.go`（+1）、`orders.go`（+1）
- `v2/rest/client.go` — 註冊新 service 到 Client struct

**原因/備註：** v1 加前綴讓呼叫端明確區分版本；v2 endpoint 覆蓋率從 57% 提升至約 95%

---

## 2026-04-01 — Add new v2 REST endpoint services and extend existing ones

**改動摘要：** 新增 6 個 REST service 檔案（account, movements, alerts, settings, liquidations, rankings），並擴充 5 個現有檔案（positions, funding, derivatives, status, orders）加入新的 API endpoint 方法

**修改的檔案：**
- `v2/rest/account.go` — 新增 AccountService（UserInfo, Summary, MarginInfo, AvailableBalance, LoginsHistory, AuditHistory）
- `v2/rest/movements.go` — 新增 MovementsService（History, AllHistory）
- `v2/rest/alerts.go` — 新增 AlertService（All, Set, Delete）
- `v2/rest/settings.go` — 新增 SettingsService（Read, Write, Delete）
- `v2/rest/liquidations.go` — 新增 LiquidationsService（History）
- `v2/rest/rankings.go` — 新增 RankingsService（History）
- `v2/rest/positions.go` — 新增 History, Audit, Snapshot, Increase, IncreaseInfo 方法
- `v2/rest/funding.go` — 新增 CancelAllOffers, Info, AutoRenew, Close 方法
- `v2/rest/derivatives.go` — 新增 CollateralLimits 方法
- `v2/rest/status.go` — 新增 DerivativeStatusHistory 方法
- `v2/rest/orders.go` — 新增 AllActive 方法
- `v2/rest/client.go` — 註冊新的 service 到 Client struct 並初始化

---

## 2026-04-01 — Fix double-dot syntax errors from broken sed operation

**改動摘要：** 修復所有 v1 test 檔案及相關檔案中因 sed 操作產生的 `..` (double-dot) 語法錯誤，改為正確的 `.` member access；同時修復被破壞的 test function 名稱（如 `TestHistor.V1Balance` -> `TestHistoryV1Balance`）

**修改的檔案：**
- `v1/trades_test.go` — `NewClient()..Trades` -> `NewClient().Trades`
- `v1/positions_test.go` — `client..Positions` -> `client.Positions`, `NewClient()..Positions` -> `NewClient().Positions`
- `v1/orders_test.go` — `NewClient()..Orders` -> `NewClient().Orders` (3 處)
- `v1/history_test.go` — `NewClient()..History` -> `NewClient().History`, 修復函數名 `TestHistor.V1Balance` / `TestHistor.V1Movements`
- `v1/ticker_test.go` — `NewClient()..Ticker` -> `NewClient().Ticker`
- `v1/deposit_test.go` — `NewClient()..Deposit` -> `NewClient().Deposit`
- `v1/lendbook_test.go` — `NewClient()..Lendbook` -> `NewClient().Lendbook` (2 處)
- `v1/doc.go` — `api..Pairs` -> `api.Pairs`
- `v1/wallet_test.go` — `NewClient()..Wallet` -> `NewClient().Wallet`, 修復函數名 `Tes.V1WithdrawCrypto` / `Tes.V1WithdrawWire`
- `v1/offers_test.go` — `NewClient()..Offers` -> `NewClient().Offers` (2 處)
- `v1/pairs_test.go` — `NewClient()..Pairs` -> `NewClient().Pairs`, 修復函數名 `TestPair.V1AllDetailed`
- `v1/stats_test.go` — `NewClient()..Stats` -> `NewClient().Stats`
- `v1/order_book_test.go` — `NewClient()..OrderBook` -> `NewClient().OrderBook`
- `tests/integration/v1/balances_test.go` — `client..Balances` -> `client.Balances`
- `tests/integration/v1/order_test.go` — `client..Orders` -> `client.Orders`
- `tests/integration/v1/order_book_test.go` — `client..OrderBook` -> `client.OrderBook`

**原因/備註：** 先前的 sed 批次重命名操作產生了 `..` 語法錯誤，共修復 27 處

---

## 2026-04-01 — Bitfinex REST v2 API endpoint inventory

**改動摘要：** 從 Bitfinex 官方文件抓取並整理所有 REST v2 API endpoint 清單（public + authenticated）

**修改的檔案：**
- 無程式碼改動，僅為資料蒐集任務

**原因/備註：** 用於比對目前 Go client 已實作的 endpoint 與官方文件的差異

---

## 2026-04-01 16:04 — P1/P2/P3 全面修復與升級

**改動摘要：** 修復 P1 靜默錯誤、升級 Go 版本、清理 deprecated API、加入 error wrapping、架構改善、新增測試

**修改的檔案：**

### P1: 錯誤傳遞修復
- `v2/rest/client.go:157` — `sign()` 錯誤從 `return "", nil` 改為 `return "", err`
- `v1/websocket.go:284-296` — auth 失敗時透過 channel 回報錯誤並關閉 ws，不再靜默 return
- `v1/websocket.go:131` — `sendSubscribeMessages()` 中 `json.Marshal` 錯誤不再忽略

### P2: Go 版本升級 & deprecated API 清理
- `go.mod` — Go 1.15 → Go 1.21
- `v1/client.go` — `ioutil.ReadAll` → `io.ReadAll`
- `v2/rest/client.go` — `ioutil.ReadAll` → `io.ReadAll`，移除 `io/ioutil` import
- 16 個 test 檔案 — `ioutil.NopCloser` → `io.NopCloser`

### P2: Error wrapping（`%w`）
- `v2/rest/transport.go` — 所有 error return 加上 `fmt.Errorf` context wrapping
- `v2/rest/client.go` — `sign()`、`NewAuthenticatedRequestWithData()` 加上 wrapping
- `v2/websocket/client.go` — `sign()` 加上 wrapping
- `v2/websocket/transport.go` — `Connect()` 加上 wrapping

### P3: 架構改善
- `v1/websocket.go` — 新增 `wsReadTimeout`(90s)、`wsWriteTimeout`(10s) 常數；`Subscribe()`、`sendSubscribeMessages()`、`ConnectPrivate()` 加上 read/write deadline；`Close()` 加上 nil 保護；`handleDataMessage()` 加上空 slice 長度檢查
- `v2/websocket/subscriptions.go` — `control()` 從 `time.Sleep` 改為 `time.Ticker` + `select`（修復 Close 時最多等 7.5s 的問題）；加入 `sync.WaitGroup` 追蹤 goroutine 生命週期；`Close()` 呼叫 `wg.Wait()` 確保 goroutine 已退出

### 新增測試
- `v2/rest/transport_test.go` — error check 順序、error wrapping、invalid URL
- `v2/rest/client_sign_test.go` — sign 函數正確性
- `v2/websocket/client_close_test.go` — double-close 不 panic、有 socket 的 close、空 socket close
- `v2/websocket/subscriptions_test.go` — control goroutine 停止、add/lookup、reset
- `v1/websocket_safety_test.go` — 惡意 payload 不 panic、valid payload、unmapped channel、nil connection close

**原因/備註：** 全面提升 codebase 穩健性，遵循 Go idiom（error wrapping、context timeout、goroutine lifecycle management）

---

## 2026-04-01 — Fix P0 crash bugs (5 issues)

**改動摘要：** 修復 5 個會導致生產環境 panic 的 P0 等級 bug

**修改的檔案：**
- `v2/rest/transport.go` — 修正 error check 順序，避免 `http.NewRequest` 失敗時對 nil `httpReq` 操作導致 nil pointer dereference
- `v1/websocket.go:321-348` — 將 heartbeat 的 `return` 改為 `continue`，避免收到第一個 heartbeat 就中斷整個 WebSocket 迴圈
- `v1/websocket.go:184-224` — 所有 type assertion 改為 comma-ok 語法，避免 API 回傳格式異常時 panic；同時檢查 chanMap key 是否存在再 send
- `v2/websocket/transport.go:178-186` — 用 `select` + `kill` channel 保護 downstream send，避免 shutdown 期間 send on closed channel
- `v2/websocket/client.go:245-264` — 用 `sync.Once` 包裝 `Close()` 防止 double-close panic；對 `c.sockets` 讀取加 mutex 保護

**原因/備註：** 這些都是靜態分析可發現的 crash bug，在網路異常、API 格式變更、或並發關閉時會觸發 panic

---

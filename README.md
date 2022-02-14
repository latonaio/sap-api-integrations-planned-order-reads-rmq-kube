# sap-api-integrations-planned-order-reads-rmq-kube  
sap-api-integrations-planned-order-reads-rmq-kube は、外部システム(特にエッジコンピューティング環境)をSAPと統合することを目的に、SAP API で 計画手配 データを取得するマイクロサービスです。  
sap-api-integrations-planned-order-reads-rmq-kube には、サンプルのAPI Json フォーマットが含まれています。  
sap-api-integrations-planned-order-reads-rmq-kube は、オンプレミス版である（＝クラウド版ではない）SAPS4HANA API の利用を前提としています。クラウド版APIを利用する場合は、ご注意ください。  
https://api.sap.com/api/OP_PLANNEDORDER_0001/overview  

## 動作環境
sap-api-integrations-planned-order-reads-rmq-kube は、主にエッジコンピューティング環境における動作にフォーカスしています。   
使用する際は、事前に下記の通り エッジコンピューティングの動作環境（推奨/必須）を用意してください。   
・ エッジ Kubernetes （推奨）    
・ AION のリソース （推奨)    
・ OS: LinuxOS （必須）    
・ CPU: ARM/AMD/Intel（いずれか必須） 

## クラウド環境での利用  
sap-api-integrations-planned-order-reads-rmq-kube は、外部システムがクラウド環境である場合にSAPと統合するときにおいても、利用可能なように設計されています。  

## RabbitMQ からの JSON Input

sap-api-integrations-planned-order-reads-rmq-kube は、Inputとして、RabbitMQ からのメッセージをJSON形式で受け取ります。 
Input の サンプルJSON は、Inputs フォルダ内にあります。  

## RabbitMQ からのメッセージ受信による イベントドリヴン の ランタイム実行

sap-api-integrations-planned-order-reads-rmq-kube は、RabbitMQ からのメッセージを受け取ると、イベントドリヴンでランタイムを実行します。  
AION の仕様では、Kubernetes 上 の 当該マイクロサービスPod は 立ち上がったまま待機状態で当該メッセージを受け取り、（コンテナ起動などの段取時間をカットして）即座にランタイムを実行します。　

## RabbitMQ への JSON Output

sap-api-integrations-planned-order-reads-rmq-kube は、Outputとして、RabbitMQ へのメッセージをJSON形式で出力します。  
Output の サンプルJSON は、Outputs フォルダ内にあります。  

## RabbitMQ の マスタサーバ環境

sap-api-integrations-planned-order-reads-rmq-kube が利用する RabbitMQ のマスタサーバ環境は、[rabbitmq-on-kubernetes](https://github.com/latonaio/rabbitmq-on-kubernetes) です。  
当該マスタサーバ環境は、同じエッジコンピューティングデバイスに配置されても、別の物理(仮想)サーバ内に配置されても、どちらでも構いません。

## RabbitMQ の Golang Runtime ライブラリ
sap-api-integrations-planned-order-reads-rmq-kube は、RabbitMQ の Golang Runtime ライブラリ として、[rabbitmq-golang-client](https://github.com/latonaio/rabbitmq-golang-client)を利用しています。

## デプロイ・稼働
sap-api-integrations-planned-order-reads-rmq-kube の デプロイ・稼働 を行うためには、aion-service-definitions の services.yml に、本レポジトリの services.yml を設定する必要があります。

kubectl apply - f 等で Deployment作成後、以下のコマンドで Pod が正しく生成されていることを確認してください。
```
$ kubectl get pods
```

## 本レポジトリ が 対応する API サービス
sap-api-integrations-planned-order-reads-rmq-kube が対応する APIサービス は、次のものです。

* APIサービス概要説明 URL: https://api.sap.com/api/OP_PLANNEDORDER_0001/overview  
* APIサービス名(=baseURL): api_plannedorder/srvd_a2x/sap/plannedorder/0001

## 本レポジトリ に 含まれる API名
sap-api-integrations-planned-order-reads-rmq-kube には、次の API をコールするためのリソースが含まれています。  

* PlannedOrderHeader（計画手配 - ヘッダ）
* PlannedOrderComponent（計画手配 - 構成品目）

## API への 値入力条件 の 初期値
sap-api-integrations-planned-order-reads-rmq-kube において、API への値入力条件の初期値は、入力ファイルレイアウトの種別毎に、次の通りとなっています。  

## SDC レイアウト

* inoutSDC.PlannedOrder.PlannedOrder（計画手配）
* inoutSDC.PlannedOrder.Material（品目）
* inoutSDC.PlannedOrder.MRPPlant（MRPプラント）
* inoutSDC.PlannedOrder.Component.Plant（プラント）

## SAP API Bussiness Hub の API の選択的コール

Latona および AION の SAP 関連リソースでは、Inputs フォルダ下の sample.json の accepter に取得したいデータの種別（＝APIの種別）を入力し、指定することができます。  
なお、同 accepter にAll(もしくは空白)の値を入力することで、全データ（＝全APIの種別）をまとめて取得することができます。  

* sample.jsonの記載例(1)  

accepter において 下記の例のように、データの種別（＝APIの種別）を指定します。  
ここでは、"Header" が指定されています。    
  
```
	"api_schema": "PlannedOrderHeader",
	"accepter": ["Header"],
	"planned_order": "39",
	"deleted": false
```
  
* 全データを取得する際のsample.jsonの記載例(2)  

全データを取得する場合、sample.json は以下のように記載します。  

```
	"api_schema": "PlannedOrderHeader",
	"accepter": ["All"],
	"planned_order": "39",
	"deleted": false
```
## 指定されたデータ種別のコール

accepter における データ種別 の指定に基づいて SAP_API_Caller 内の caller.go で API がコールされます。  
caller.go の func() 毎 の 以下の箇所が、指定された API をコールするソースコードです。  

```
func (c *SAPAPICaller) AsyncGetPlannedOrder(plannedOrder, material, mRPPlant, plant string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "Header":
			func() {
				c.Header(plannedOrder)
				wg.Done()
			}()
		case "HeaderMaterialPlant":
			func() {
				c.HeaderMaterialPlant(material, mRPPlant)
				wg.Done()
			}()
		case "ComponentMaterialPlant":
			func() {
				c.ComponentMaterialPlant(material, plant)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}
```

## SAP API Business Hub における API サービス の バージョン と バージョン におけるデータレイアウトの相違

SAP API Business Hub における API サービス のうちの 殆どの API サービス のBASE URLのフォーマットは、"API_(リポジトリ名)_SRV" であり、殆どの API サービス 間 の データレイアウトは統一されています。   
従って、Latona および AION における リソースにおいても、データレイアウトが統一されています。    
一方、本レポジトリ に関わる API である Planned Order のサービスは、BASE URLのフォーマットが他のAPIサービスと異なります。      
その結果、本レポジトリ内の一部のAPIのデータレイアウトが、他のAPIサービスのものと異なっています。  

#### BASE URLが "API_(リポジトリ名)_SRV" のフォーマットである API サービス の データレイアウト（=responses）  
BASE URLが "API_{リポジトリ名}_SRV" のフォーマットであるAPIサービスのデータレイアウト（=responses）は、例えば、次の通りです。  
```
type ToProductionOrderItem struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID   string `json:"id"`
				URI  string `json:"uri"`
				Type string `json:"type"`
			} `json:"__metadata"`
			ManufacturingOrder             string      `json:"ManufacturingOrder"`
			ManufacturingOrderItem         string      `json:"ManufacturingOrderItem"`
			ManufacturingOrderCategory     string      `json:"ManufacturingOrderCategory"`
			ManufacturingOrderType         string      `json:"ManufacturingOrderType"`
			IsCompletelyDelivered          bool        `json:"IsCompletelyDelivered"`
			Material                       string      `json:"Material"`
			ProductionPlant                string      `json:"ProductionPlant"`
			Plant                          string      `json:"Plant"`
			MRPArea                        string      `json:"MRPArea"`
			QuantityDistributionKey        string      `json:"QuantityDistributionKey"`
			MaterialGoodsReceiptDuration   string      `json:"MaterialGoodsReceiptDuration"`
			StorageLocation                string      `json:"StorageLocation"`
			Batch                          string      `json:"Batch"`
			InventoryUsabilityCode         string      `json:"InventoryUsabilityCode"`
			GoodsRecipientName             string      `json:"GoodsRecipientName"`
			UnloadingPointName             string      `json:"UnloadingPointName"`
			MfgOrderItemPlndDeliveryDate   string      `json:"MfgOrderItemPlndDeliveryDate"`
			MfgOrderItemActualDeliveryDate string      `json:"MfgOrderItemActualDeliveryDate"`
			ProductionUnit                 string      `json:"ProductionUnit"`
			MfgOrderItemPlannedTotalQty    string      `json:"MfgOrderItemPlannedTotalQty"`
			MfgOrderItemPlannedScrapQty    string      `json:"MfgOrderItemPlannedScrapQty"`
			MfgOrderItemGoodsReceiptQty    string      `json:"MfgOrderItemGoodsReceiptQty"`
			MfgOrderItemActualDeviationQty string      `json:"MfgOrderItemActualDeviationQty"`
		} `json:"results"`
	} `json:"d"`
}

```

#### BASE URL が "api_plannedorder/srvd_a2x/sap/plannedorder/0001" である Planned Order の APIサービス の データレイアウト（=responses）  
BASE URL が "api_plannedorder/srvd_a2x/sap/plannedorder/0001" である Planned Order の APIサービス の データレイアウト（=responses）は、例えば、次の通りです。  

```
type Header struct {
	Value             []struct {
		PlannedOrder                   string        `json:"PlannedOrder"`
		PlannedOrderType               string        `json:"PlannedOrderType"`
		PlannedOrderProfile            string        `json:"PlannedOrderProfile"`
		Material                       string        `json:"Material"`
		MaterialName                   string        `json:"MaterialName"`
		ProductionPlant                string        `json:"ProductionPlant"`
		MRPPlant                       string        `json:"MRPPlant"`
		MRPArea                        string        `json:"MRPArea"`
		ProductionVersion              string        `json:"ProductionVersion"`
		MaterialProcurementCategory    string        `json:"MaterialProcurementCategory"`
		MaterialProcurementType        string        `json:"MaterialProcurementType"`
		StorageLocation                string        `json:"StorageLocation"`
		BaseUnit                       string        `json:"BaseUnit"`
		TotalQuantity                  int           `json:"TotalQuantity"`
		PlndOrderPlannedScrapQty       int           `json:"PlndOrderPlannedScrapQty"`
		GoodsReceiptQty                int           `json:"GoodsReceiptQty"`
		IssuedQuantity                 int           `json:"IssuedQuantity"`
		PlndOrderPlannedStartDate      string        `json:"PlndOrderPlannedStartDate"`
		PlndOrderPlannedStartTime      string        `json:"PlndOrderPlannedStartTime"`
		PlndOrderPlannedEndDate        string        `json:"PlndOrderPlannedEndDate"`
		PlndOrderPlannedEndTime        string        `json:"PlndOrderPlannedEndTime"`
		PlannedOrderOpeningDate        string        `json:"PlannedOrderOpeningDate"`
		PlannedOrderLastChangeDateTime string        `json:"PlannedOrderLastChangeDateTime"`
		ProductionStartDate            string        `json:"ProductionStartDate"`
		ProductionEndDate              string        `json:"ProductionEndDate"`
		SalesOrder                     string        `json:"SalesOrder"`
		SalesOrderItem                 string        `json:"SalesOrderItem"`
		Customer                       string        `json:"Customer"`
		WBSElementInternalID           string        `json:"WBSElementInternalID"`
		WBSElementExternalID           string        `json:"WBSElementExternalID"`
		WBSDescription                 string        `json:"WBSDescription"`
		AccountAssignmentCategory      string        `json:"AccountAssignmentCategory"`
		Reservation                    string        `json:"Reservation"`
		PlannedOrderLongText           string        `json:"PlannedOrderLongText"`
		MRPController                  string        `json:"MRPController"`
		ProductionSupervisor           string        `json:"ProductionSupervisor"`
		PurchasingGroup                string        `json:"PurchasingGroup"`
		PurchasingOrganization         string        `json:"PurchasingOrganization"`
		FixedSupplier                  string        `json:"FixedSupplier"`
		PurchasingDocument             string        `json:"PurchasingDocument"`
		PurchasingDocumentItem         string        `json:"PurchasingDocumentItem"`
		SupplierName                   string        `json:"SupplierName"`
		PlannedOrderIsFirm             bool          `json:"PlannedOrderIsFirm"`
		PlannedOrderIsConvertible      bool          `json:"PlannedOrderIsConvertible"`
		PlannedOrderBOMIsFixed         bool          `json:"PlannedOrderBOMIsFixed"`
		PlannedOrderCapacityIsDsptchd  bool          `json:"PlannedOrderCapacityIsDsptchd"`
		CapacityRequirement            string        `json:"CapacityRequirement"`
		BillOfOperationsVariant        string        `json:"BillOfOperationsVariant"`
		CapacityRequirementOrigin      string        `json:"CapacityRequirementOrigin"`
		BillOfOperationsType           string        `json:"BillOfOperationsType"`
		BillOfOperationsGroup          string        `json:"BillOfOperationsGroup"`
		LastScheduledDate              string        `json:"LastScheduledDate"`
		ScheduledBasicEndDate          string        `json:"ScheduledBasicEndDate"`
		ScheduledBasicEndTime          string        `json:"ScheduledBasicEndTime"`
		ScheduledBasicStartDate        string        `json:"ScheduledBasicStartDate"`
		ScheduledBasicStartTime        string        `json:"ScheduledBasicStartTime"`
		SchedulingType                 string        `json:"SchedulingType"`
	} `json:"value"`
}

```
このように、BASE URLが "API_(リポジトリ名)_SRV" のフォーマットである API サービス の データレイアウトと、 Planned Order の データレイアウトは、D、Results、Metadata、Value の配列構造を持っているか持っていないかという点が異なります。  

## Output  
本マイクロサービスでは、[golang-logging-library-for-sap](https://github.com/latonaio/golang-logging-library-for-sap) により、以下のようなデータがJSON形式で出力されます。  
以下の sample.json の例は、SAP 計画手配 の ヘッダデータ が取得された結果の JSON の例です。  
以下の項目のうち、"PlannedOrder" ～ "SchedulingType" は、/SAP_API_Output_Formatter/type.go 内 の Type Header {} による出力結果です。"cursor" ～ "time"は、golang-logging-library-for-sap による 定型フォーマットの出力結果です。  

```
{
	"cursor": "/Users/latona2/bitbucket/sap-api-integrations-planned-order-reads/SAP_API_Caller/caller.go#L63",
	"function": "sap-api-integrations-planned-order-reads/SAP_API_Caller.(*SAPAPICaller).Header",
	"level": "INFO",
	"message": [
		{
			"PlannedOrder": "39",
			"PlannedOrderType": "PE",
			"PlannedOrderProfile": "",
			"Material": "FG233",
			"MaterialName": "FERT 233, PD, Repetitive Manuf.",
			"ProductionPlant": "1010",
			"MRPPlant": "1010",
			"MRPArea": "1010",
			"ProductionVersion": "0001",
			"MaterialProcurementCategory": "E",
			"MaterialProcurementType": "E",
			"StorageLocation": "101A",
			"BaseUnit": "PC",
			"TotalQuantity": 20,
			"PlndOrderPlannedScrapQty": 0,
			"GoodsReceiptQty": 0,
			"IssuedQuantity": 0,
			"PlndOrderPlannedStartDate": "2017-03-21",
			"PlndOrderPlannedStartTime": "00:00:00",
			"PlndOrderPlannedEndDate": "2017-03-21",
			"PlndOrderPlannedEndTime": "00:00:00",
			"PlannedOrderOpeningDate": "2017-03-07",
			"PlannedOrderLastChangeDateTime": "2017-08-30T04:58:36Z",
			"ProductionStartDate": "",
			"ProductionEndDate": "",
			"SalesOrder": "",
			"SalesOrderItem": "0",
			"Customer": "",
			"WBSElementInternalID": "0",
			"WBSElementExternalID": "",
			"WBSDescription": "",
			"AccountAssignmentCategory": "",
			"Reservation": "1691",
			"PlannedOrderLongText": "",
			"MRPController": "001",
			"ProductionSupervisor": "YB1",
			"PurchasingGroup": "",
			"PurchasingOrganization": "",
			"FixedSupplier": "",
			"PurchasingDocument": "",
			"PurchasingDocumentItem": "0",
			"SupplierName": "",
			"PlannedOrderIsFirm": true,
			"PlannedOrderIsConvertible": false,
			"PlannedOrderBOMIsFixed": false,
			"PlannedOrderCapacityIsDsptchd": false,
			"CapacityRequirement": "10000005168",
			"BillOfOperationsVariant": "1",
			"CapacityRequirementOrigin": "2",
			"BillOfOperationsType": "R",
			"BillOfOperationsGroup": "41010009",
			"LastScheduledDate": "2017-08-30",
			"ScheduledBasicEndDate": "2017-09-01",
			"ScheduledBasicEndTime": "07:02:15",
			"ScheduledBasicStartDate": "2017-09-01",
			"ScheduledBasicStartTime": "07:00:00",
			"SchedulingType": "2"
		}
	],
	"time": "2022-01-28T15:02:22+09:00"
}
```

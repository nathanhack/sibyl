// Code generated by ogen, DO NOT EDIT.

package ogent

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AddTicker implements addTicker operation.
//
// Queue for adding entities by ticker.
//
// POST /rest/entities/add/{ticker}
func (UnimplementedHandler) AddTicker(ctx context.Context, params AddTickerParams) error {
	return ht.ErrNotImplemented
}

// CreateBarGroup implements createBarGroup operation.
//
// Creates a new BarGroup and persists it to storage.
//
// POST /rest/bar-groups
func (UnimplementedHandler) CreateBarGroup(ctx context.Context, req *CreateBarGroupReq) (r CreateBarGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateBarRecord implements createBarRecord operation.
//
// Creates a new BarRecord and persists it to storage.
//
// POST /rest/bar-records
func (UnimplementedHandler) CreateBarRecord(ctx context.Context, req *CreateBarRecordReq) (r CreateBarRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateBarTimeRange implements createBarTimeRange operation.
//
// Creates a new BarTimeRange and persists it to storage.
//
// POST /rest/bar-time-ranges
func (UnimplementedHandler) CreateBarTimeRange(ctx context.Context, req *CreateBarTimeRangeReq) (r CreateBarTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateDataSource implements createDataSource operation.
//
// Creates a new DataSource and persists it to storage.
//
// POST /rest/data-sources
func (UnimplementedHandler) CreateDataSource(ctx context.Context, req *CreateDataSourceReq) (r CreateDataSourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateDividend implements createDividend operation.
//
// Creates a new Dividend and persists it to storage.
//
// POST /rest/dividends
func (UnimplementedHandler) CreateDividend(ctx context.Context, req *CreateDividendReq) (r CreateDividendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateEntity implements createEntity operation.
//
// Creates a new Entity and persists it to storage.
//
// POST /rest/entities
func (UnimplementedHandler) CreateEntity(ctx context.Context, req *CreateEntityReq) (r CreateEntityRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateExchange implements createExchange operation.
//
// Creates a new Exchange and persists it to storage.
//
// POST /rest/exchanges
func (UnimplementedHandler) CreateExchange(ctx context.Context, req *CreateExchangeReq) (r CreateExchangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateFinancial implements createFinancial operation.
//
// Creates a new Financial and persists it to storage.
//
// POST /rest/financials
func (UnimplementedHandler) CreateFinancial(ctx context.Context, req *CreateFinancialReq) (r CreateFinancialRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateInterval implements createInterval operation.
//
// Creates a new Interval and persists it to storage.
//
// POST /rest/intervals
func (UnimplementedHandler) CreateInterval(ctx context.Context, req *CreateIntervalReq) (r CreateIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateMarketHours implements createMarketHours operation.
//
// Creates a new MarketHours and persists it to storage.
//
// POST /rest/market-hours
func (UnimplementedHandler) CreateMarketHours(ctx context.Context, req *CreateMarketHoursReq) (r CreateMarketHoursRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateMarketInfo implements createMarketInfo operation.
//
// Creates a new MarketInfo and persists it to storage.
//
// POST /rest/market-infos
func (UnimplementedHandler) CreateMarketInfo(ctx context.Context, req *CreateMarketInfoReq) (r CreateMarketInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateSplit implements createSplit operation.
//
// Creates a new Split and persists it to storage.
//
// POST /rest/splits
func (UnimplementedHandler) CreateSplit(ctx context.Context, req *CreateSplitReq) (r CreateSplitRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateTradeCondition implements createTradeCondition operation.
//
// Creates a new TradeCondition and persists it to storage.
//
// POST /rest/trade-conditions
func (UnimplementedHandler) CreateTradeCondition(ctx context.Context, req *CreateTradeConditionReq) (r CreateTradeConditionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateTradeCorrection implements createTradeCorrection operation.
//
// Creates a new TradeCorrection and persists it to storage.
//
// POST /rest/trade-corrections
func (UnimplementedHandler) CreateTradeCorrection(ctx context.Context, req *CreateTradeCorrectionReq) (r CreateTradeCorrectionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateTradeRecord implements createTradeRecord operation.
//
// Creates a new TradeRecord and persists it to storage.
//
// POST /rest/trade-records
func (UnimplementedHandler) CreateTradeRecord(ctx context.Context, req *CreateTradeRecordReq) (r CreateTradeRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateTradeTimeRange implements createTradeTimeRange operation.
//
// Creates a new TradeTimeRange and persists it to storage.
//
// POST /rest/trade-time-ranges
func (UnimplementedHandler) CreateTradeTimeRange(ctx context.Context, req *CreateTradeTimeRangeReq) (r CreateTradeTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteBarGroup implements deleteBarGroup operation.
//
// Deletes the BarGroup with the requested ID.
//
// DELETE /rest/bar-groups/{id}
func (UnimplementedHandler) DeleteBarGroup(ctx context.Context, params DeleteBarGroupParams) (r DeleteBarGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteBarRecord implements deleteBarRecord operation.
//
// Deletes the BarRecord with the requested ID.
//
// DELETE /rest/bar-records/{id}
func (UnimplementedHandler) DeleteBarRecord(ctx context.Context, params DeleteBarRecordParams) (r DeleteBarRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteBarTimeRange implements deleteBarTimeRange operation.
//
// Deletes the BarTimeRange with the requested ID.
//
// DELETE /rest/bar-time-ranges/{id}
func (UnimplementedHandler) DeleteBarTimeRange(ctx context.Context, params DeleteBarTimeRangeParams) (r DeleteBarTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteDataSource implements deleteDataSource operation.
//
// Deletes the DataSource with the requested ID.
//
// DELETE /rest/data-sources/{id}
func (UnimplementedHandler) DeleteDataSource(ctx context.Context, params DeleteDataSourceParams) (r DeleteDataSourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteDividend implements deleteDividend operation.
//
// Deletes the Dividend with the requested ID.
//
// DELETE /rest/dividends/{id}
func (UnimplementedHandler) DeleteDividend(ctx context.Context, params DeleteDividendParams) (r DeleteDividendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteEntity implements deleteEntity operation.
//
// Deletes the Entity with the requested ID.
//
// DELETE /rest/entities/{id}
func (UnimplementedHandler) DeleteEntity(ctx context.Context, params DeleteEntityParams) (r DeleteEntityRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteExchange implements deleteExchange operation.
//
// Deletes the Exchange with the requested ID.
//
// DELETE /rest/exchanges/{id}
func (UnimplementedHandler) DeleteExchange(ctx context.Context, params DeleteExchangeParams) (r DeleteExchangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteFinancial implements deleteFinancial operation.
//
// Deletes the Financial with the requested ID.
//
// DELETE /rest/financials/{id}
func (UnimplementedHandler) DeleteFinancial(ctx context.Context, params DeleteFinancialParams) (r DeleteFinancialRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteInterval implements deleteInterval operation.
//
// Deletes the Interval with the requested ID.
//
// DELETE /rest/intervals/{id}
func (UnimplementedHandler) DeleteInterval(ctx context.Context, params DeleteIntervalParams) (r DeleteIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteMarketHours implements deleteMarketHours operation.
//
// Deletes the MarketHours with the requested ID.
//
// DELETE /rest/market-hours/{id}
func (UnimplementedHandler) DeleteMarketHours(ctx context.Context, params DeleteMarketHoursParams) (r DeleteMarketHoursRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteMarketInfo implements deleteMarketInfo operation.
//
// Deletes the MarketInfo with the requested ID.
//
// DELETE /rest/market-infos/{id}
func (UnimplementedHandler) DeleteMarketInfo(ctx context.Context, params DeleteMarketInfoParams) (r DeleteMarketInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteSplit implements deleteSplit operation.
//
// Deletes the Split with the requested ID.
//
// DELETE /rest/splits/{id}
func (UnimplementedHandler) DeleteSplit(ctx context.Context, params DeleteSplitParams) (r DeleteSplitRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteTradeCondition implements deleteTradeCondition operation.
//
// Deletes the TradeCondition with the requested ID.
//
// DELETE /rest/trade-conditions/{id}
func (UnimplementedHandler) DeleteTradeCondition(ctx context.Context, params DeleteTradeConditionParams) (r DeleteTradeConditionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteTradeCorrection implements deleteTradeCorrection operation.
//
// Deletes the TradeCorrection with the requested ID.
//
// DELETE /rest/trade-corrections/{id}
func (UnimplementedHandler) DeleteTradeCorrection(ctx context.Context, params DeleteTradeCorrectionParams) (r DeleteTradeCorrectionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteTradeRecord implements deleteTradeRecord operation.
//
// Deletes the TradeRecord with the requested ID.
//
// DELETE /rest/trade-records/{id}
func (UnimplementedHandler) DeleteTradeRecord(ctx context.Context, params DeleteTradeRecordParams) (r DeleteTradeRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteTradeTimeRange implements deleteTradeTimeRange operation.
//
// Deletes the TradeTimeRange with the requested ID.
//
// DELETE /rest/trade-time-ranges/{id}
func (UnimplementedHandler) DeleteTradeTimeRange(ctx context.Context, params DeleteTradeTimeRangeParams) (r DeleteTradeTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListBarGroup implements listBarGroup operation.
//
// List BarGroups.
//
// GET /rest/bar-groups
func (UnimplementedHandler) ListBarGroup(ctx context.Context, params ListBarGroupParams) (r ListBarGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListBarGroupRecords implements listBarGroupRecords operation.
//
// List attached Records.
//
// GET /rest/bar-groups/{id}/records
func (UnimplementedHandler) ListBarGroupRecords(ctx context.Context, params ListBarGroupRecordsParams) (r ListBarGroupRecordsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListBarRecord implements listBarRecord operation.
//
// List BarRecords.
//
// GET /rest/bar-records
func (UnimplementedHandler) ListBarRecord(ctx context.Context, params ListBarRecordParams) (r ListBarRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListBarTimeRange implements listBarTimeRange operation.
//
// List BarTimeRanges.
//
// GET /rest/bar-time-ranges
func (UnimplementedHandler) ListBarTimeRange(ctx context.Context, params ListBarTimeRangeParams) (r ListBarTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListBarTimeRangeGroups implements listBarTimeRangeGroups operation.
//
// List attached Groups.
//
// GET /rest/bar-time-ranges/{id}/groups
func (UnimplementedHandler) ListBarTimeRangeGroups(ctx context.Context, params ListBarTimeRangeGroupsParams) (r ListBarTimeRangeGroupsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListDataSource implements listDataSource operation.
//
// List DataSources.
//
// GET /rest/data-sources
func (UnimplementedHandler) ListDataSource(ctx context.Context, params ListDataSourceParams) (r ListDataSourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListDataSourceIntervals implements listDataSourceIntervals operation.
//
// List attached Intervals.
//
// GET /rest/data-sources/{id}/intervals
func (UnimplementedHandler) ListDataSourceIntervals(ctx context.Context, params ListDataSourceIntervalsParams) (r ListDataSourceIntervalsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListDividend implements listDividend operation.
//
// List Dividends.
//
// GET /rest/dividends
func (UnimplementedHandler) ListDividend(ctx context.Context, params ListDividendParams) (r ListDividendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListDividendStock implements listDividendStock operation.
//
// List attached Stocks.
//
// GET /rest/dividends/{id}/stock
func (UnimplementedHandler) ListDividendStock(ctx context.Context, params ListDividendStockParams) (r ListDividendStockRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListEntity implements listEntity operation.
//
// List Entities.
//
// GET /rest/entities
func (UnimplementedHandler) ListEntity(ctx context.Context, params ListEntityParams) (r ListEntityRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListEntityDividends implements listEntityDividends operation.
//
// List attached Dividends.
//
// GET /rest/entities/{id}/dividends
func (UnimplementedHandler) ListEntityDividends(ctx context.Context, params ListEntityDividendsParams) (r ListEntityDividendsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListEntityExchanges implements listEntityExchanges operation.
//
// List attached Exchanges.
//
// GET /rest/entities/{id}/exchanges
func (UnimplementedHandler) ListEntityExchanges(ctx context.Context, params ListEntityExchangesParams) (r ListEntityExchangesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListEntityFinancials implements listEntityFinancials operation.
//
// List attached Financials.
//
// GET /rest/entities/{id}/financials
func (UnimplementedHandler) ListEntityFinancials(ctx context.Context, params ListEntityFinancialsParams) (r ListEntityFinancialsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListEntityIntervals implements listEntityIntervals operation.
//
// List attached Intervals.
//
// GET /rest/entities/{id}/intervals
func (UnimplementedHandler) ListEntityIntervals(ctx context.Context, params ListEntityIntervalsParams) (r ListEntityIntervalsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListEntitySplits implements listEntitySplits operation.
//
// List attached Splits.
//
// GET /rest/entities/{id}/splits
func (UnimplementedHandler) ListEntitySplits(ctx context.Context, params ListEntitySplitsParams) (r ListEntitySplitsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListExchange implements listExchange operation.
//
// List Exchanges.
//
// GET /rest/exchanges
func (UnimplementedHandler) ListExchange(ctx context.Context, params ListExchangeParams) (r ListExchangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListExchangeStocks implements listExchangeStocks operation.
//
// List attached Stocks.
//
// GET /rest/exchanges/{id}/stocks
func (UnimplementedHandler) ListExchangeStocks(ctx context.Context, params ListExchangeStocksParams) (r ListExchangeStocksRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListFinancial implements listFinancial operation.
//
// List Financials.
//
// GET /rest/financials
func (UnimplementedHandler) ListFinancial(ctx context.Context, params ListFinancialParams) (r ListFinancialRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListFinancialStock implements listFinancialStock operation.
//
// List attached Stocks.
//
// GET /rest/financials/{id}/stock
func (UnimplementedHandler) ListFinancialStock(ctx context.Context, params ListFinancialStockParams) (r ListFinancialStockRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListInterval implements listInterval operation.
//
// List Intervals.
//
// GET /rest/intervals
func (UnimplementedHandler) ListInterval(ctx context.Context, params ListIntervalParams) (r ListIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListIntervalBars implements listIntervalBars operation.
//
// List attached Bars.
//
// GET /rest/intervals/{id}/bars
func (UnimplementedHandler) ListIntervalBars(ctx context.Context, params ListIntervalBarsParams) (r ListIntervalBarsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListIntervalTrades implements listIntervalTrades operation.
//
// List attached Trades.
//
// GET /rest/intervals/{id}/trades
func (UnimplementedHandler) ListIntervalTrades(ctx context.Context, params ListIntervalTradesParams) (r ListIntervalTradesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListMarketHours implements listMarketHours operation.
//
// List MarketHours.
//
// GET /rest/market-hours
func (UnimplementedHandler) ListMarketHours(ctx context.Context, params ListMarketHoursParams) (r ListMarketHoursRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListMarketInfo implements listMarketInfo operation.
//
// List MarketInfos.
//
// GET /rest/market-infos
func (UnimplementedHandler) ListMarketInfo(ctx context.Context, params ListMarketInfoParams) (r ListMarketInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListMarketInfoHours implements listMarketInfoHours operation.
//
// List attached Hours.
//
// GET /rest/market-infos/{id}/hours
func (UnimplementedHandler) ListMarketInfoHours(ctx context.Context, params ListMarketInfoHoursParams) (r ListMarketInfoHoursRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListSplit implements listSplit operation.
//
// List Splits.
//
// GET /rest/splits
func (UnimplementedHandler) ListSplit(ctx context.Context, params ListSplitParams) (r ListSplitRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeCondition implements listTradeCondition operation.
//
// List TradeConditions.
//
// GET /rest/trade-conditions
func (UnimplementedHandler) ListTradeCondition(ctx context.Context, params ListTradeConditionParams) (r ListTradeConditionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeConditionRecord implements listTradeConditionRecord operation.
//
// List attached Records.
//
// GET /rest/trade-conditions/{id}/record
func (UnimplementedHandler) ListTradeConditionRecord(ctx context.Context, params ListTradeConditionRecordParams) (r ListTradeConditionRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeCorrection implements listTradeCorrection operation.
//
// List TradeCorrections.
//
// GET /rest/trade-corrections
func (UnimplementedHandler) ListTradeCorrection(ctx context.Context, params ListTradeCorrectionParams) (r ListTradeCorrectionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeCorrectionRecord implements listTradeCorrectionRecord operation.
//
// List attached Records.
//
// GET /rest/trade-corrections/{id}/record
func (UnimplementedHandler) ListTradeCorrectionRecord(ctx context.Context, params ListTradeCorrectionRecordParams) (r ListTradeCorrectionRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeRecord implements listTradeRecord operation.
//
// List TradeRecords.
//
// GET /rest/trade-records
func (UnimplementedHandler) ListTradeRecord(ctx context.Context, params ListTradeRecordParams) (r ListTradeRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeRecordConditions implements listTradeRecordConditions operation.
//
// List attached Conditions.
//
// GET /rest/trade-records/{id}/conditions
func (UnimplementedHandler) ListTradeRecordConditions(ctx context.Context, params ListTradeRecordConditionsParams) (r ListTradeRecordConditionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeRecordCorrection implements listTradeRecordCorrection operation.
//
// List attached Corrections.
//
// GET /rest/trade-records/{id}/correction
func (UnimplementedHandler) ListTradeRecordCorrection(ctx context.Context, params ListTradeRecordCorrectionParams) (r ListTradeRecordCorrectionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeRecordExchange implements listTradeRecordExchange operation.
//
// List attached Exchanges.
//
// GET /rest/trade-records/{id}/exchange
func (UnimplementedHandler) ListTradeRecordExchange(ctx context.Context, params ListTradeRecordExchangeParams) (r ListTradeRecordExchangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeTimeRange implements listTradeTimeRange operation.
//
// List TradeTimeRanges.
//
// GET /rest/trade-time-ranges
func (UnimplementedHandler) ListTradeTimeRange(ctx context.Context, params ListTradeTimeRangeParams) (r ListTradeTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTradeTimeRangeRecords implements listTradeTimeRangeRecords operation.
//
// List attached Records.
//
// GET /rest/trade-time-ranges/{id}/records
func (UnimplementedHandler) ListTradeTimeRangeRecords(ctx context.Context, params ListTradeTimeRangeRecordsParams) (r ListTradeTimeRangeRecordsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadBarGroup implements readBarGroup operation.
//
// Finds the BarGroup with the requested ID and returns it.
//
// GET /rest/bar-groups/{id}
func (UnimplementedHandler) ReadBarGroup(ctx context.Context, params ReadBarGroupParams) (r ReadBarGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadBarGroupTimeRange implements readBarGroupTimeRange operation.
//
// Find the attached BarTimeRange of the BarGroup with the given ID.
//
// GET /rest/bar-groups/{id}/time-range
func (UnimplementedHandler) ReadBarGroupTimeRange(ctx context.Context, params ReadBarGroupTimeRangeParams) (r ReadBarGroupTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadBarRecord implements readBarRecord operation.
//
// Finds the BarRecord with the requested ID and returns it.
//
// GET /rest/bar-records/{id}
func (UnimplementedHandler) ReadBarRecord(ctx context.Context, params ReadBarRecordParams) (r ReadBarRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadBarRecordGroup implements readBarRecordGroup operation.
//
// Find the attached BarGroup of the BarRecord with the given ID.
//
// GET /rest/bar-records/{id}/group
func (UnimplementedHandler) ReadBarRecordGroup(ctx context.Context, params ReadBarRecordGroupParams) (r ReadBarRecordGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadBarTimeRange implements readBarTimeRange operation.
//
// Finds the BarTimeRange with the requested ID and returns it.
//
// GET /rest/bar-time-ranges/{id}
func (UnimplementedHandler) ReadBarTimeRange(ctx context.Context, params ReadBarTimeRangeParams) (r ReadBarTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadBarTimeRangeInterval implements readBarTimeRangeInterval operation.
//
// Find the attached Interval of the BarTimeRange with the given ID.
//
// GET /rest/bar-time-ranges/{id}/interval
func (UnimplementedHandler) ReadBarTimeRangeInterval(ctx context.Context, params ReadBarTimeRangeIntervalParams) (r ReadBarTimeRangeIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadDataSource implements readDataSource operation.
//
// Finds the DataSource with the requested ID and returns it.
//
// GET /rest/data-sources/{id}
func (UnimplementedHandler) ReadDataSource(ctx context.Context, params ReadDataSourceParams) (r ReadDataSourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadDividend implements readDividend operation.
//
// Finds the Dividend with the requested ID and returns it.
//
// GET /rest/dividends/{id}
func (UnimplementedHandler) ReadDividend(ctx context.Context, params ReadDividendParams) (r ReadDividendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadEntity implements readEntity operation.
//
// Finds the Entity with the requested ID and returns it.
//
// GET /rest/entities/{id}
func (UnimplementedHandler) ReadEntity(ctx context.Context, params ReadEntityParams) (r ReadEntityRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadExchange implements readExchange operation.
//
// Finds the Exchange with the requested ID and returns it.
//
// GET /rest/exchanges/{id}
func (UnimplementedHandler) ReadExchange(ctx context.Context, params ReadExchangeParams) (r ReadExchangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadFinancial implements readFinancial operation.
//
// Finds the Financial with the requested ID and returns it.
//
// GET /rest/financials/{id}
func (UnimplementedHandler) ReadFinancial(ctx context.Context, params ReadFinancialParams) (r ReadFinancialRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadInterval implements readInterval operation.
//
// Finds the Interval with the requested ID and returns it.
//
// GET /rest/intervals/{id}
func (UnimplementedHandler) ReadInterval(ctx context.Context, params ReadIntervalParams) (r ReadIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadIntervalDataSource implements readIntervalDataSource operation.
//
// Find the attached DataSource of the Interval with the given ID.
//
// GET /rest/intervals/{id}/data-source
func (UnimplementedHandler) ReadIntervalDataSource(ctx context.Context, params ReadIntervalDataSourceParams) (r ReadIntervalDataSourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadIntervalStock implements readIntervalStock operation.
//
// Find the attached Entity of the Interval with the given ID.
//
// GET /rest/intervals/{id}/stock
func (UnimplementedHandler) ReadIntervalStock(ctx context.Context, params ReadIntervalStockParams) (r ReadIntervalStockRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadMarketHours implements readMarketHours operation.
//
// Finds the MarketHours with the requested ID and returns it.
//
// GET /rest/market-hours/{id}
func (UnimplementedHandler) ReadMarketHours(ctx context.Context, params ReadMarketHoursParams) (r ReadMarketHoursRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadMarketHoursMarketInfo implements readMarketHoursMarketInfo operation.
//
// Find the attached MarketInfo of the MarketHours with the given ID.
//
// GET /rest/market-hours/{id}/market-info
func (UnimplementedHandler) ReadMarketHoursMarketInfo(ctx context.Context, params ReadMarketHoursMarketInfoParams) (r ReadMarketHoursMarketInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadMarketInfo implements readMarketInfo operation.
//
// Finds the MarketInfo with the requested ID and returns it.
//
// GET /rest/market-infos/{id}
func (UnimplementedHandler) ReadMarketInfo(ctx context.Context, params ReadMarketInfoParams) (r ReadMarketInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadSplit implements readSplit operation.
//
// Finds the Split with the requested ID and returns it.
//
// GET /rest/splits/{id}
func (UnimplementedHandler) ReadSplit(ctx context.Context, params ReadSplitParams) (r ReadSplitRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadSplitStock implements readSplitStock operation.
//
// Find the attached Entity of the Split with the given ID.
//
// GET /rest/splits/{id}/stock
func (UnimplementedHandler) ReadSplitStock(ctx context.Context, params ReadSplitStockParams) (r ReadSplitStockRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadTradeCondition implements readTradeCondition operation.
//
// Finds the TradeCondition with the requested ID and returns it.
//
// GET /rest/trade-conditions/{id}
func (UnimplementedHandler) ReadTradeCondition(ctx context.Context, params ReadTradeConditionParams) (r ReadTradeConditionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadTradeCorrection implements readTradeCorrection operation.
//
// Finds the TradeCorrection with the requested ID and returns it.
//
// GET /rest/trade-corrections/{id}
func (UnimplementedHandler) ReadTradeCorrection(ctx context.Context, params ReadTradeCorrectionParams) (r ReadTradeCorrectionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadTradeRecord implements readTradeRecord operation.
//
// Finds the TradeRecord with the requested ID and returns it.
//
// GET /rest/trade-records/{id}
func (UnimplementedHandler) ReadTradeRecord(ctx context.Context, params ReadTradeRecordParams) (r ReadTradeRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadTradeRecordTimeRange implements readTradeRecordTimeRange operation.
//
// Find the attached TradeTimeRange of the TradeRecord with the given ID.
//
// GET /rest/trade-records/{id}/time-range
func (UnimplementedHandler) ReadTradeRecordTimeRange(ctx context.Context, params ReadTradeRecordTimeRangeParams) (r ReadTradeRecordTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadTradeTimeRange implements readTradeTimeRange operation.
//
// Finds the TradeTimeRange with the requested ID and returns it.
//
// GET /rest/trade-time-ranges/{id}
func (UnimplementedHandler) ReadTradeTimeRange(ctx context.Context, params ReadTradeTimeRangeParams) (r ReadTradeTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadTradeTimeRangeInterval implements readTradeTimeRangeInterval operation.
//
// Find the attached Interval of the TradeTimeRange with the given ID.
//
// GET /rest/trade-time-ranges/{id}/interval
func (UnimplementedHandler) ReadTradeTimeRangeInterval(ctx context.Context, params ReadTradeTimeRangeIntervalParams) (r ReadTradeTimeRangeIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SearchTicker implements searchTicker operation.
//
// Searches for entities by ticker.
//
// GET /rest/search/{ticker}
func (UnimplementedHandler) SearchTicker(ctx context.Context, params SearchTickerParams) (r *SearchTickerOK, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateBarGroup implements updateBarGroup operation.
//
// Updates a BarGroup and persists changes to storage.
//
// PATCH /rest/bar-groups/{id}
func (UnimplementedHandler) UpdateBarGroup(ctx context.Context, req *UpdateBarGroupReq, params UpdateBarGroupParams) (r UpdateBarGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateBarRecord implements updateBarRecord operation.
//
// Updates a BarRecord and persists changes to storage.
//
// PATCH /rest/bar-records/{id}
func (UnimplementedHandler) UpdateBarRecord(ctx context.Context, req *UpdateBarRecordReq, params UpdateBarRecordParams) (r UpdateBarRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateBarTimeRange implements updateBarTimeRange operation.
//
// Updates a BarTimeRange and persists changes to storage.
//
// PATCH /rest/bar-time-ranges/{id}
func (UnimplementedHandler) UpdateBarTimeRange(ctx context.Context, req *UpdateBarTimeRangeReq, params UpdateBarTimeRangeParams) (r UpdateBarTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateDataSource implements updateDataSource operation.
//
// Updates a DataSource and persists changes to storage.
//
// PATCH /rest/data-sources/{id}
func (UnimplementedHandler) UpdateDataSource(ctx context.Context, req *UpdateDataSourceReq, params UpdateDataSourceParams) (r UpdateDataSourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateDividend implements updateDividend operation.
//
// Updates a Dividend and persists changes to storage.
//
// PATCH /rest/dividends/{id}
func (UnimplementedHandler) UpdateDividend(ctx context.Context, req *UpdateDividendReq, params UpdateDividendParams) (r UpdateDividendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateEntity implements updateEntity operation.
//
// Updates a Entity and persists changes to storage.
//
// PATCH /rest/entities/{id}
func (UnimplementedHandler) UpdateEntity(ctx context.Context, req *UpdateEntityReq, params UpdateEntityParams) (r UpdateEntityRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateExchange implements updateExchange operation.
//
// Updates a Exchange and persists changes to storage.
//
// PATCH /rest/exchanges/{id}
func (UnimplementedHandler) UpdateExchange(ctx context.Context, req *UpdateExchangeReq, params UpdateExchangeParams) (r UpdateExchangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateFinancial implements updateFinancial operation.
//
// Updates a Financial and persists changes to storage.
//
// PATCH /rest/financials/{id}
func (UnimplementedHandler) UpdateFinancial(ctx context.Context, req *UpdateFinancialReq, params UpdateFinancialParams) (r UpdateFinancialRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateInterval implements updateInterval operation.
//
// Updates a Interval and persists changes to storage.
//
// PATCH /rest/intervals/{id}
func (UnimplementedHandler) UpdateInterval(ctx context.Context, req *UpdateIntervalReq, params UpdateIntervalParams) (r UpdateIntervalRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateMarketHours implements updateMarketHours operation.
//
// Updates a MarketHours and persists changes to storage.
//
// PATCH /rest/market-hours/{id}
func (UnimplementedHandler) UpdateMarketHours(ctx context.Context, req *UpdateMarketHoursReq, params UpdateMarketHoursParams) (r UpdateMarketHoursRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateMarketInfo implements updateMarketInfo operation.
//
// Updates a MarketInfo and persists changes to storage.
//
// PATCH /rest/market-infos/{id}
func (UnimplementedHandler) UpdateMarketInfo(ctx context.Context, req *UpdateMarketInfoReq, params UpdateMarketInfoParams) (r UpdateMarketInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateSplit implements updateSplit operation.
//
// Updates a Split and persists changes to storage.
//
// PATCH /rest/splits/{id}
func (UnimplementedHandler) UpdateSplit(ctx context.Context, req *UpdateSplitReq, params UpdateSplitParams) (r UpdateSplitRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateTradeCondition implements updateTradeCondition operation.
//
// Updates a TradeCondition and persists changes to storage.
//
// PATCH /rest/trade-conditions/{id}
func (UnimplementedHandler) UpdateTradeCondition(ctx context.Context, req *UpdateTradeConditionReq, params UpdateTradeConditionParams) (r UpdateTradeConditionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateTradeCorrection implements updateTradeCorrection operation.
//
// Updates a TradeCorrection and persists changes to storage.
//
// PATCH /rest/trade-corrections/{id}
func (UnimplementedHandler) UpdateTradeCorrection(ctx context.Context, req *UpdateTradeCorrectionReq, params UpdateTradeCorrectionParams) (r UpdateTradeCorrectionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateTradeRecord implements updateTradeRecord operation.
//
// Updates a TradeRecord and persists changes to storage.
//
// PATCH /rest/trade-records/{id}
func (UnimplementedHandler) UpdateTradeRecord(ctx context.Context, req *UpdateTradeRecordReq, params UpdateTradeRecordParams) (r UpdateTradeRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateTradeTimeRange implements updateTradeTimeRange operation.
//
// Updates a TradeTimeRange and persists changes to storage.
//
// PATCH /rest/trade-time-ranges/{id}
func (UnimplementedHandler) UpdateTradeTimeRange(ctx context.Context, req *UpdateTradeTimeRangeReq, params UpdateTradeTimeRangeParams) (r UpdateTradeTimeRangeRes, _ error) {
	return r, ht.ErrNotImplemented
}

package currency

import "strings"

type CurrencyConf struct {
	Currency  string
	Label     string
	Symbol    string
	Pairs     []string
	Pools     []string
	Explorers ExplorerConf
	Unit      string
}

type ExplorerConf struct {
	BaseUri        string
	AddressUri     string
	TransactionUri string
}

type CurrencyConfigMapping string

const (
	CurrencyLabelMap    CurrencyConfigMapping = "pub:map:currency:label"
	CurrencySymbolMap   CurrencyConfigMapping = "pub:map:currency:sym"
	CurrencyUnitMap     CurrencyConfigMapping = "pub:map:currency:unit"
	CurrencyExplorerMap CurrencyConfigMapping = "pub:map:currency:explorer"
	CurrencyExchangeMap CurrencyConfigMapping = "pub:list:pair:exchange"
)

type RawCurrencyConf struct {
	Mapping string
	Data    interface{}
}

func parseCurrencyLabelMap(config map[string]CurrencyConf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Label = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := CurrencyConf{}
			cfg.Label = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseCurrencySymbMap(config map[string]CurrencyConf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Symbol = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := CurrencyConf{}
			cfg.Symbol = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseCurrencyUnitMap(config map[string]CurrencyConf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Unit = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := CurrencyConf{}
			cfg.Unit = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseCurrencyExplorerMap(config map[string]CurrencyConf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		explorers := data[1].([]interface{})
		var cfg CurrencyConf
		if val, ok := config[cur]; ok {
			cfg = val
		} else {
			// create new empty config instance
			cc := CurrencyConf{}
			cc.Currency = cur
			cfg = cc
		}
		ec := ExplorerConf{
			explorers[0].(string),
			explorers[1].(string),
			explorers[2].(string),
		}
		cfg.Explorers = ec
		config[cur] = cfg
	}
}

func parseCurrencyExchangeMap(config map[string]CurrencyConf, raw []interface{}) {
	for _, rs := range raw {
		symbol := rs.(string)
		var base, quote string

		if len(symbol) > 6 {
			base = strings.Split(symbol, ":")[0]
			quote = strings.Split(symbol, ":")[1]
		} else {
			base = symbol[3:]
			quote = symbol[:3]
		}

		// append if base exists in configs
		if val, ok := config[base]; ok {
			val.Pairs = append(val.Pairs, symbol)
			config[base] = val
		}

		// append if quote exists in configs
		if val, ok := config[quote]; ok {
			val.Pairs = append(val.Pairs, symbol)
			config[quote] = val
		}
	}
}

func NewCurrencyConfFromRaw(raw []RawCurrencyConf) ([]CurrencyConf, error) {
	configMap := make(map[string]CurrencyConf)
	for _, r := range raw {
		switch CurrencyConfigMapping(r.Mapping) {
		case CurrencyLabelMap:
			data := r.Data.([]interface{})
			parseCurrencyLabelMap(configMap, data)
		case CurrencySymbolMap:
			data := r.Data.([]interface{})
			parseCurrencySymbMap(configMap, data)
		case CurrencyUnitMap:
			data := r.Data.([]interface{})
			parseCurrencyUnitMap(configMap, data)
		case CurrencyExplorerMap:
			data := r.Data.([]interface{})
			parseCurrencyExplorerMap(configMap, data)
		case CurrencyExchangeMap:
			data := r.Data.([]interface{})
			parseCurrencyExchangeMap(configMap, data)
		}
	}

	// convert map to array
	configs := make([]CurrencyConf, 0)
	for _, v := range configMap {
		configs = append(configs, v)
	}

	return configs, nil
}

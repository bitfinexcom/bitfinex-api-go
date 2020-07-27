package currency

import "strings"

type Conf struct {
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

type ConfigMapping string

const (
	LabelMap    ConfigMapping = "pub:map:currency:label"
	SymbolMap   ConfigMapping = "pub:map:currency:sym"
	UnitMap     ConfigMapping = "pub:map:currency:unit"
	ExplorerMap ConfigMapping = "pub:map:currency:explorer"
	ExchangeMap ConfigMapping = "pub:list:pair:exchange"
)

type RawConf struct {
	Mapping string
	Data    interface{}
}

func parseLabelMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Label = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := Conf{}
			cfg.Label = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseSymbMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Symbol = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := Conf{}
			cfg.Symbol = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseUnitMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		if val, ok := config[cur]; ok {
			// add value
			val.Unit = data[1].(string)
			config[cur] = val
		} else {
			// create new empty config instance
			cfg := Conf{}
			cfg.Unit = data[1].(string)
			cfg.Currency = cur
			config[cur] = cfg
		}
	}
}

func parseExplorerMap(config map[string]Conf, raw []interface{}) {
	for _, rawLabel := range raw {
		data := rawLabel.([]interface{})
		cur := data[0].(string)
		explorers := data[1].([]interface{})
		var cfg Conf
		if val, ok := config[cur]; ok {
			cfg = val
		} else {
			// create new empty config instance
			cc := Conf{}
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

func parseExchangeMap(config map[string]Conf, raw []interface{}) {
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

func FromRaw(raw []RawConf) ([]Conf, error) {
	configMap := make(map[string]Conf)
	for _, r := range raw {
		switch ConfigMapping(r.Mapping) {
		case LabelMap:
			data := r.Data.([]interface{})
			parseLabelMap(configMap, data)
		case SymbolMap:
			data := r.Data.([]interface{})
			parseSymbMap(configMap, data)
		case UnitMap:
			data := r.Data.([]interface{})
			parseUnitMap(configMap, data)
		case ExplorerMap:
			data := r.Data.([]interface{})
			parseExplorerMap(configMap, data)
		case ExchangeMap:
			data := r.Data.([]interface{})
			parseExchangeMap(configMap, data)
		}
	}

	// convert map to array
	configs := make([]Conf, 0)
	for _, v := range configMap {
		configs = append(configs, v)
	}

	return configs, nil
}

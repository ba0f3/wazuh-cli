package cmd

import (
	"fmt"
	"strings"

	"github.com/ba0f3/wazuh-cli/internal/config"
	"github.com/ba0f3/wazuh-cli/internal/indexer"
	"github.com/ba0f3/wazuh-cli/internal/output"
	"github.com/spf13/cobra"
)

var globalIndexer *indexer.Client

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Query Wazuh alerts from the Indexer (OpenSearch)",
	Long: `Query Wazuh alerts directly from the Wazuh Indexer (OpenSearch).
	
Alerts are not stored in the Wazuh Manager API, so this command uses a separate
connection to the Indexer (typically port 9200). You must configure the indexer
URL before using this command:

  wazuh-cli config set indexer_url https://indexer-node:9200
  wazuh-cli config set indexer_user admin
  wazuh-cli config set indexer_password secret
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Resolve configuration
		cfg, err := config.Load(configPath, &flagOverrides)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		globalCfg = cfg
		globalFmt = output.New(cfg.Output, cfg.Pretty)

		if cfg.IndexerURL == "" {
			return fmt.Errorf("indexer_url is not configured. Hint: run 'wazuh-cli config set indexer_url https://...'")
		}

		c, err := indexer.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("creating indexer client: %w", err)
		}
		globalIndexer = c

		return nil
	},
}

var alertListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent alerts",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		level, _ := cmd.Flags().GetInt("level")
		agentID, _ := cmd.Flags().GetString("agent-id")
		agentName, _ := cmd.Flags().GetString("agent-name")
		ruleID, _ := cmd.Flags().GetString("rule-id")
		ruleGroup, _ := cmd.Flags().GetString("rule-group")
		queryStr, _ := cmd.Flags().GetString("query")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		sortField, _ := cmd.Flags().GetString("sort")
		indexPattern, _ := cmd.Flags().GetString("index")

		if indexPattern == "" {
			indexPattern = globalCfg.IndexerIndex
			if indexPattern == "" {
				indexPattern = "wazuh-alerts-4.x-*"
			}
		}

		queryDSL := buildSearchQuery(limit, level, agentID, agentName, ruleID, ruleGroup, queryStr, from, to, sortField)

		resp, err := globalIndexer.Search(indexPattern, queryDSL)
		if err != nil {
			return err
		}

		// Extract hits to return just the source documents for output
		var docs []interface{}
		for _, hit := range resp.Hits.Hits {
			docs = append(docs, hit.Source)
		}

		if len(docs) == 0 {
			mustWrite([]interface{}{})
			return nil
		}

		mustWrite(docs)
		return nil
	},
}

var alertGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a specific alert by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		indexPattern, _ := cmd.Flags().GetString("index")

		if indexPattern == "" {
			indexPattern = globalCfg.IndexerIndex
			if indexPattern == "" {
				indexPattern = "wazuh-alerts-4.x-*"
			}
		}

		hit, err := globalIndexer.Get(indexPattern, id)
		if err != nil {
			return err
		}

		mustWrite(hit.Source)
		return nil
	},
}

var alertStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get aggregated alert statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupBy, _ := cmd.Flags().GetString("group-by")
		level, _ := cmd.Flags().GetInt("level")
		agentID, _ := cmd.Flags().GetString("agent-id")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		indexPattern, _ := cmd.Flags().GetString("index")

		if indexPattern == "" {
			indexPattern = globalCfg.IndexerIndex
			if indexPattern == "" {
				indexPattern = "wazuh-alerts-4.x-*"
			}
		}

		queryDSL := buildStatsQuery(groupBy, level, agentID, from, to)

		resp, err := globalIndexer.Search(indexPattern, queryDSL)
		if err != nil {
			return err
		}

		if agg, ok := resp.Aggregations["stats"]; ok {
			mustWrite(agg.Buckets)
		} else {
			mustWrite([]interface{}{})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(alertCmd)
	alertCmd.AddCommand(alertListCmd, alertGetCmd, alertStatsCmd)

	// List flags
	alertListCmd.Flags().Int("limit", 50, "maximum number of alerts to return")
	alertListCmd.Flags().Int("level", 0, "minimum rule level")
	alertListCmd.Flags().String("agent-id", "", "filter by agent ID")
	alertListCmd.Flags().String("agent-name", "", "filter by agent name")
	alertListCmd.Flags().String("rule-id", "", "filter by rule ID")
	alertListCmd.Flags().String("rule-group", "", "filter by rule group")
	alertListCmd.Flags().String("query", "", "raw Lucene query string")
	alertListCmd.Flags().String("from", "now-24h", "start time (e.g. now-1h, 2024-01-01T00:00:00Z)")
	alertListCmd.Flags().String("to", "now", "end time")
	alertListCmd.Flags().String("sort", "timestamp:desc", "sort field and direction")
	alertListCmd.Flags().String("index", "", "override index pattern (default: wazuh-alerts-4.x-*)")

	// Get flags
	alertGetCmd.Flags().String("index", "", "override index pattern (default: wazuh-alerts-4.x-*)")

	// Stats flags
	alertStatsCmd.Flags().String("group-by", "level", "field to group by: level, agent, rule")
	alertStatsCmd.Flags().Int("level", 0, "minimum rule level")
	alertStatsCmd.Flags().String("agent-id", "", "filter by agent ID")
	alertStatsCmd.Flags().String("from", "now-24h", "start time")
	alertStatsCmd.Flags().String("to", "now", "end time")
	alertStatsCmd.Flags().String("index", "", "override index pattern (default: wazuh-alerts-4.x-*)")
}

// buildSearchQuery constructs an OpenSearch DSL query for alert list.
func buildSearchQuery(limit, level int, agentID, agentName, ruleID, ruleGroup, queryStr, from, to, sortField string) map[string]interface{} {
	must := []map[string]interface{}{}

	if level > 0 {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"rule.level": map[string]interface{}{"gte": level},
			},
		})
	}
	if agentID != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"agent.id": agentID},
		})
	}
	if agentName != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{"agent.name": agentName},
		})
	}
	if ruleID != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"rule.id": ruleID},
		})
	}
	if ruleGroup != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{"rule.groups": ruleGroup},
		})
	}
	if queryStr != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{"query": queryStr},
		})
	}

	// Time range
	must = append(must, map[string]interface{}{
		"range": map[string]interface{}{
			"timestamp": map[string]interface{}{
				"gte": from,
				"lte": to,
			},
		},
	})

	sortMap := map[string]interface{}{}
	parts := strings.Split(sortField, ":")
	if len(parts) == 2 {
		sortMap[parts[0]] = parts[1]
	} else {
		sortMap[sortField] = "asc"
	}

	query := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": []map[string]interface{}{
			sortMap,
		},
	}

	return query
}

// buildStatsQuery constructs an OpenSearch DSL query for alert aggregations.
func buildStatsQuery(groupBy string, level int, agentID, from, to string) map[string]interface{} {
	must := []map[string]interface{}{}

	if level > 0 {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"rule.level": map[string]interface{}{"gte": level},
			},
		})
	}
	if agentID != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"agent.id": agentID},
		})
	}

	must = append(must, map[string]interface{}{
		"range": map[string]interface{}{
			"timestamp": map[string]interface{}{
				"gte": from,
				"lte": to,
			},
		},
	})

	fieldMap := map[string]string{
		"level": "rule.level",
		"agent": "agent.name",
		"rule":  "rule.id",
	}

	aggField := "rule.level"
	if f, ok := fieldMap[groupBy]; ok {
		aggField = f
	} else {
		aggField = groupBy
	}

	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"aggs": map[string]interface{}{
			"stats": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": aggField,
					"size":  50,
				},
			},
		},
	}

	return query
}

package json

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/pryingbytez/pryingdeep/configs"
	"github.com/pryingbytez/pryingdeep/models"
	"github.com/pryingbytez/pryingdeep/pkg/exporters"
	"github.com/pryingbytez/pryingdeep/pkg/querybuilder"
)

var JSONCmd = &cobra.Command{
	Use:   "json",
	Short: "Export the crawled data to json",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExportJSON(cmd, args)
	},
}

var (
	rawSQL       = false
	rawFilePath  = "pkg/querybuilder/queries/select.sql"
	criteria     map[string]string
	associations = "all"
	sortBy       = "status_code"
	sortOrder    = "asc"
	limit        = 0
	filePath     = "data.json"
)

func init() {
	JSONCmd.Flags().BoolVarP(&rawSQL, "raw", "r", rawSQL, "--raw to use raw sql queries that you provide. All other flags except silent, rawFilePath and filepath will not matter.")
	JSONCmd.Flags().StringVarP(&rawFilePath, "raw-sql-filepath", "p", rawFilePath, "-rp to specify the file path to the sql file. Only use this flag if you specify -raw")
	JSONCmd.Flags().StringToStringVarP(&criteria, "criteria", "q", criteria, "JSON-formatted criteria, e.g., -q 'title=test,\"url=LIKE example.com\"'")
	JSONCmd.Flags().StringVarP(&associations, "associations", "a", associations, "-a WP,E,P,C")
	JSONCmd.Flags().StringVarP(&sortBy, "sort-by", "b", sortBy, "SortBy e.g -> -b title")
	JSONCmd.Flags().StringVarP(&sortOrder, "sort-order", "o", sortOrder, "SortOrder e.g -> -o ASC || -b DESC. Only use this flag if you use SortBy")
	JSONCmd.Flags().IntVarP(&limit, "limit", "l", limit, "Limit e.g -> -l 100 -> 100 items will be taken from the database. Default limit will acquire all results from the database")
	JSONCmd.Flags().StringVarP(&filePath, "filepath", "f", filePath, "FilePath -f myfilepath")
	JSONCmd.MarkFlagsRequiredTogether("raw", "raw-sql-filepath")

	cli := configs.NewCLIConfig()
	JSONCmd.Flags().VisitAll(cli.ConfigureViper("exporter"))
}

func ExportJSON(cmd *cobra.Command, args []string) error {
	var data interface{}
	var err error

	db := models.GetDB()
	exporterConfig := configs.LoadExporterConfig()
	setExportOptions(&exporterConfig)

	if !rawSQL {
		color.HiMagenta("[+] Constructing query...")
		qb := querybuilder.NewQueryBuilder(
			exporterConfig.Criteria,
			exporterConfig.Associations,
			exporterConfig.SortBy,
			exporterConfig.SortOrder,
			exporterConfig.Limit,
		)
		data = qb.ConstructQuery(db)
	} else {
		color.HiRed("[+] Reading raw query...")
		qb := querybuilder.NewQueryBuilder(nil, "", "", "", 0)
		err, data = qb.Raw(db, rawFilePath)
		if err != nil {
			return err
		}
	}
	err = exporters.ExportDataToJSON(data, exporterConfig.Filepath)
	if err != nil {
		return err
	}
	return nil
}

func setExportOptions(eC *configs.Exporter) {
	if len(criteria) != 0 {
		eC.Criteria = make(map[string]interface{})
		for key, value := range criteria {
			eC.Criteria[key] = value
		}
	}
	if associations != "" {
		eC.Associations = associations
	}
	if sortBy != "" {
		eC.SortBy = sortBy
	}

	if sortOrder != "" {
		eC.SortOrder = sortOrder
	}

	if limit >= 0 {
		eC.Limit = limit
	}

	if filePath != "" {
		eC.Filepath = filePath
	}

}

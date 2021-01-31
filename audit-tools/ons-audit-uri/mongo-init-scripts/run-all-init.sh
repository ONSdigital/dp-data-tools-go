# stop on error
set -e


echo "initialising: home_page"
mongo home_page-init.js

echo "initialising: taxonomy_landing_page"
mongo taxonomy_landing_page-init.js

echo "initialising: product_page"
mongo product_page-init.js

echo "initialising: article"
mongo article-init.js
echo "initialising: article_download"
mongo article_download-init.js
echo "initialising: bulletin"
mongo bulletin-init.js
echo "initialising: compendium_data"
mongo compendium_data-init.js
echo "initialising: compendium_landing_page"
mongo compendium_landing_page-init.js
echo "initialising: dataset_landing_page"
mongo dataset_landing_page-init.js
echo "initialising: static_methodology"
mongo static_methodology-init.js
echo "initialising: static_methodology_download"
mongo static_methodology_download-init.js
echo "initialising: static_qmi"
mongo static_qmi-init.js
echo "initialising: timeseries"
mongo timeseries-init.js

echo "initialising: chart"
mongo chart-init.js

echo "initialising: table"
mongo table-init.js
echo "initialising: equation"
mongo equation-init.js
echo "initialising: image"
mongo image-init.js
echo "initialising: release"
mongo release-init.js
echo "initialising: list"
mongo list-init.js
echo "initialising: static_page"
mongo static_page-init.js
echo "initialising: static_adhoc"
mongo static_adhoc-init.js
echo "initialising: reference_tables"
mongo reference_tables-init.js
echo "initialising: compendium_chapter"
mongo compendium_chapter-init.js
echo "initialising: static_landing_page"
mongo static_landing_page-init.js
echo "initialising: static_article"
mongo static_article-init.js
echo "initialising: dataset"
mongo dataset-init.js
echo "initialising: timeseries_dataset"
mongo timeseries_dataset-init.js

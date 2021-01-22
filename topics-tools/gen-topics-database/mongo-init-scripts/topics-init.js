db = db.getSiblingDB('topics')
db.topics.remove({})
db.topics.insertOne({
    "id": "topic_root",
    "next": {
        "id": "topic_root",
        "description": "root page",
        "title": "The root page",
        "keywords": [
            "root 1",
            "root 2"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/topic_root",
                "id": "topic_root"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/topic_root/subtopics"
            }
        },
        "subtopics_ids": [
            "businessindustryandtrade",
            "economy",
            "employmentandlabourmarket",
            "peoplepopulationandcommunity"
        ]
    },
    "current": {
        "id": "topic_root",
        "description": "root page",
        "title": "The root page",
        "keywords": [
            "root 1",
            "root 2"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/topic_root",
                "id": "topic_root"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/topic_root/subtopics"
            }
        },
        "subtopics_ids": [
            "businessindustryandtrade",
            "economy",
            "employmentandlabourmarket",
            "peoplepopulationandcommunity"
        ]
    }
})
db.topics.insertOne({
    "id": "businessindustryandtrade",
    "next": {
        "id": "businessindustryandtrade",
        "description": "Activities of businesses and industry in the UK, including data on the production and trade of goods and services, sales by retailers, characteristics of businesses, the construction and manufacturing sectors, and international trade.",
        "title": "Business, industry and trade",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessindustryandtrade",
                "id": "businessindustryandtrade"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/businessindustryandtrade/subtopics"
            }
        },
        "subtopics_ids": [
            "business",
            "changestobusiness",
            "constructionindustry",
            "internationaltrade",
            "itandinternetindustry",
            "manufacturingandproductionindustry",
            "retailindustry",
            "tourismindustry"
        ]
    },
    "current": {
        "id": "businessindustryandtrade",
        "description": "Activities of businesses and industry in the UK, including data on the production and trade of goods and services, sales by retailers, characteristics of businesses, the construction and manufacturing sectors, and international trade.",
        "title": "Business, industry and trade",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessindustryandtrade",
                "id": "businessindustryandtrade"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/businessindustryandtrade/subtopics"
            }
        },
        "subtopics_ids": [
            "business",
            "changestobusiness",
            "constructionindustry",
            "internationaltrade",
            "itandinternetindustry",
            "manufacturingandproductionindustry",
            "retailindustry",
            "tourismindustry"
        ]
    }
})
db.topics.insertOne({
    "id": "business",
    "next": {
        "id": "business",
        "description": "UK businesses registered for VAT and PAYE with regional breakdowns, including data on size (employment and turnover) and activity (type of industry), research and development, and business services.",
        "title": "Business",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/business",
                "id": "business"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/business/subtopics"
            }
        },
        "subtopics_ids": [
            "activitysizeandlocation",
            "businessinnovation",
            "businessservices"
        ]
    },
    "current": {
        "id": "business",
        "description": "UK businesses registered for VAT and PAYE with regional breakdowns, including data on size (employment and turnover) and activity (type of industry), research and development, and business services.",
        "title": "Business",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/business",
                "id": "business"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/business/subtopics"
            }
        },
        "subtopics_ids": [
            "activitysizeandlocation",
            "businessinnovation",
            "businessservices"
        ]
    }
})
db.topics.insertOne({
    "id": "activitysizeandlocation",
    "next": {
        "id": "activitysizeandlocation",
        "description": "UK businesses broken down by location, industry, legal status and employment, and UK enterprise levels measured by turnover.",
        "title": "Activity, size and location",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/activitysizeandlocation",
                "id": "activitysizeandlocation"
            },
            "content": {
                "href": "http://localhost:25300/topics/activitysizeandlocation/content"
            }
        }
    },
    "current": {
        "id": "activitysizeandlocation",
        "description": "UK businesses broken down by location, industry, legal status and employment, and UK enterprise levels measured by turnover.",
        "title": "Activity, size and location",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/activitysizeandlocation",
                "id": "activitysizeandlocation"
            },
            "content": {
                "href": "http://localhost:25300/topics/activitysizeandlocation/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "businessinnovation",
    "next": {
        "id": "businessinnovation",
        "description": "Research and development in the UK carried out or funded by business enterprises, higher education, government (including research councils) and private non-profit organisations.",
        "title": "Business innovation",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessinnovation",
                "id": "businessinnovation"
            },
            "content": {
                "href": "http://localhost:25300/topics/businessinnovation/content"
            }
        }
    },
    "current": {
        "id": "businessinnovation",
        "description": "Research and development in the UK carried out or funded by business enterprises, higher education, government (including research councils) and private non-profit organisations.",
        "title": "Business innovation",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessinnovation",
                "id": "businessinnovation"
            },
            "content": {
                "href": "http://localhost:25300/topics/businessinnovation/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "businessservices",
    "next": {
        "id": "businessservices",
        "description": "** no description **",
        "title": "Business services",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessservices",
                "id": "businessservices"
            },
            "content": {
                "href": "http://localhost:25300/topics/businessservices/content"
            }
        }
    },
    "current": {
        "id": "businessservices",
        "description": "** no description **",
        "title": "Business services",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessservices",
                "id": "businessservices"
            },
            "content": {
                "href": "http://localhost:25300/topics/businessservices/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "changestobusiness",
    "next": {
        "id": "changestobusiness",
        "description": "UK business growth, survival and change over time. These figures are an informal indicator of confidence in the UK economy.",
        "title": "Changes to business",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/changestobusiness",
                "id": "changestobusiness"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/changestobusiness/subtopics"
            }
        },
        "subtopics_ids": [
            "bankruptcyinsolvency",
            "businessbirthsdeathsandsurvivalrates",
            "mergersandacquisitions"
        ]
    },
    "current": {
        "id": "changestobusiness",
        "description": "UK business growth, survival and change over time. These figures are an informal indicator of confidence in the UK economy.",
        "title": "Changes to business",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/changestobusiness",
                "id": "changestobusiness"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/changestobusiness/subtopics"
            }
        },
        "subtopics_ids": [
            "bankruptcyinsolvency",
            "businessbirthsdeathsandsurvivalrates",
            "mergersandacquisitions"
        ]
    }
})
db.topics.insertOne({
    "id": "bankruptcyinsolvency",
    "next": {
        "id": "bankruptcyinsolvency",
        "description": "** no description **",
        "title": "Bankruptcy/insolvency",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/bankruptcyinsolvency",
                "id": "bankruptcyinsolvency"
            },
            "content": {
                "href": "http://localhost:25300/topics/bankruptcyinsolvency/content"
            }
        }
    },
    "current": {
        "id": "bankruptcyinsolvency",
        "description": "** no description **",
        "title": "Bankruptcy/insolvency",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/bankruptcyinsolvency",
                "id": "bankruptcyinsolvency"
            },
            "content": {
                "href": "http://localhost:25300/topics/bankruptcyinsolvency/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "businessbirthsdeathsandsurvivalrates",
    "next": {
        "id": "businessbirthsdeathsandsurvivalrates",
        "description": "Demography of UK businesses: active businesses, new registrations for VAT and PAYE (births), cessation of trading (deaths), and duration of trading (survival rates). ",
        "title": "Business births, deaths and survival rates",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessbirthsdeathsandsurvivalrates",
                "id": "businessbirthsdeathsandsurvivalrates"
            },
            "content": {
                "href": "http://localhost:25300/topics/businessbirthsdeathsandsurvivalrates/content"
            }
        }
    },
    "current": {
        "id": "businessbirthsdeathsandsurvivalrates",
        "description": "Demography of UK businesses: active businesses, new registrations for VAT and PAYE (births), cessation of trading (deaths), and duration of trading (survival rates). ",
        "title": "Business births, deaths and survival rates",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/businessbirthsdeathsandsurvivalrates",
                "id": "businessbirthsdeathsandsurvivalrates"
            },
            "content": {
                "href": "http://localhost:25300/topics/businessbirthsdeathsandsurvivalrates/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "mergersandacquisitions",
    "next": {
        "id": "mergersandacquisitions",
        "description": "Business mergers and acquisitions involving UK companies, including de-mergers and disposals, where the transaction value is £1 million or more.",
        "title": "Mergers and acquisitions",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/mergersandacquisitions",
                "id": "mergersandacquisitions"
            },
            "content": {
                "href": "http://localhost:25300/topics/mergersandacquisitions/content"
            }
        }
    },
    "current": {
        "id": "mergersandacquisitions",
        "description": "Business mergers and acquisitions involving UK companies, including de-mergers and disposals, where the transaction value is £1 million or more.",
        "title": "Mergers and acquisitions",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/mergersandacquisitions",
                "id": "mergersandacquisitions"
            },
            "content": {
                "href": "http://localhost:25300/topics/mergersandacquisitions/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "constructionindustry",
    "next": {
        "id": "constructionindustry",
        "description": "Construction of new buildings and repairs or alterations to existing properties in Great Britain measured by the amount charged for the work, including work by civil engineering companies. ",
        "title": "Construction industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/constructionindustry",
                "id": "constructionindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/constructionindustry/content"
            }
        }
    },
    "current": {
        "id": "constructionindustry",
        "description": "Construction of new buildings and repairs or alterations to existing properties in Great Britain measured by the amount charged for the work, including work by civil engineering companies. ",
        "title": "Construction industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/constructionindustry",
                "id": "constructionindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/constructionindustry/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "internationaltrade",
    "next": {
        "id": "internationaltrade",
        "description": "Trade in goods and services across the UK's international borders, including total imports and exports, the types of goods and services traded and general trends in international trade. ",
        "title": "International trade",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/internationaltrade",
                "id": "internationaltrade"
            },
            "content": {
                "href": "http://localhost:25300/topics/internationaltrade/content"
            }
        }
    },
    "current": {
        "id": "internationaltrade",
        "description": "Trade in goods and services across the UK's international borders, including total imports and exports, the types of goods and services traded and general trends in international trade. ",
        "title": "International trade",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/internationaltrade",
                "id": "internationaltrade"
            },
            "content": {
                "href": "http://localhost:25300/topics/internationaltrade/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "itandinternetindustry",
    "next": {
        "id": "itandinternetindustry",
        "description": "Internet sales by businesses in the UK (total value and as a percentage of all retail sales) and the percentage of businesses that have a website and broadband connection. These figures indicate the importance of the internet to UK businesses.",
        "title": "IT and internet industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/itandinternetindustry",
                "id": "itandinternetindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/itandinternetindustry/content"
            }
        }
    },
    "current": {
        "id": "itandinternetindustry",
        "description": "Internet sales by businesses in the UK (total value and as a percentage of all retail sales) and the percentage of businesses that have a website and broadband connection. These figures indicate the importance of the internet to UK businesses.",
        "title": "IT and internet industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/itandinternetindustry",
                "id": "itandinternetindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/itandinternetindustry/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "manufacturingandproductionindustry",
    "next": {
        "id": "manufacturingandproductionindustry",
        "description": "UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and UK manufactures' sales by product and industrial division, with EU comparisons.",
        "title": "Manufacturing and production industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/manufacturingandproductionindustry",
                "id": "manufacturingandproductionindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/manufacturingandproductionindustry/content"
            }
        }
    },
    "current": {
        "id": "manufacturingandproductionindustry",
        "description": "UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and UK manufactures' sales by product and industrial division, with EU comparisons.",
        "title": "Manufacturing and production industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/manufacturingandproductionindustry",
                "id": "manufacturingandproductionindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/manufacturingandproductionindustry/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "retailindustry",
    "next": {
        "id": "retailindustry",
        "description": "** no description **",
        "title": "Retail industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/retailindustry",
                "id": "retailindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/retailindustry/content"
            }
        }
    },
    "current": {
        "id": "retailindustry",
        "description": "** no description **",
        "title": "Retail industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/retailindustry",
                "id": "retailindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/retailindustry/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "tourismindustry",
    "next": {
        "id": "tourismindustry",
        "description": "Tourism and travel (including accommodation services, food and beverage services, passenger transport services, vehicle hire, travel agencies and sports, recreational and conference services), employment levels and output of the tourism industry, the number of visitors to the UK and the amount they spend.",
        "title": "Tourism industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/tourismindustry",
                "id": "tourismindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/tourismindustry/content"
            }
        }
    },
    "current": {
        "id": "tourismindustry",
        "description": "Tourism and travel (including accommodation services, food and beverage services, passenger transport services, vehicle hire, travel agencies and sports, recreational and conference services), employment levels and output of the tourism industry, the number of visitors to the UK and the amount they spend.",
        "title": "Tourism industry",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/tourismindustry",
                "id": "tourismindustry"
            },
            "content": {
                "href": "http://localhost:25300/topics/tourismindustry/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "economy",
    "next": {
        "id": "economy",
        "description": "UK economic activity covering production, distribution, consumption and trade of goods and services. Individuals, businesses, organisations and governments all affect the development of the economy.",
        "title": "Economy",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/economy",
                "id": "economy"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/economy/subtopics"
            }
        },
        "subtopics_ids": [
            "economicoutputandproductivity",
            "environmentalaccounts",
            "governmentpublicsectorandtaxes",
            "grossdomesticproductgdp",
            "grossvalueaddedgva",
            "inflationandpriceindices",
            "investmentspensionsandtrusts",
            "nationalaccounts",
            "regionalaccounts"
        ],
        "spotlight": [
            {
                "href": "/economy/economicoutputandproductivity/output/articles/economicactivityfasterindicatorsuk/latest",
                "title": "Research Output: Economic activity, faster indicators, UK"
            }
        ]
    },
    "current": {
        "id": "economy",
        "description": "UK economic activity covering production, distribution, consumption and trade of goods and services. Individuals, businesses, organisations and governments all affect the development of the economy.",
        "title": "Economy",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/economy",
                "id": "economy"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/economy/subtopics"
            }
        },
        "subtopics_ids": [
            "economicoutputandproductivity",
            "environmentalaccounts",
            "governmentpublicsectorandtaxes",
            "grossdomesticproductgdp",
            "grossvalueaddedgva",
            "inflationandpriceindices",
            "investmentspensionsandtrusts",
            "nationalaccounts",
            "regionalaccounts"
        ],
        "spotlight": [
            {
                "href": "/economy/economicoutputandproductivity/output/articles/economicactivityfasterindicatorsuk/latest",
                "title": "Research Output: Economic activity, faster indicators, UK"
            }
        ]
    }
})
db.topics.insertOne({
    "id": "economicoutputandproductivity",
    "next": {
        "id": "economicoutputandproductivity",
        "description": "Manufacturing, production and services indices (measuring total economic output) and productivity (measuring efficiency, expressed as a ratio of output to input over a given period of time, for example output per person per hour).",
        "title": "Economic output and productivity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/economicoutputandproductivity",
                "id": "economicoutputandproductivity"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/economicoutputandproductivity/subtopics"
            }
        },
        "subtopics_ids": [
            "output",
            "productivitymeasures",
            "publicservicesproductivity"
        ]
    },
    "current": {
        "id": "economicoutputandproductivity",
        "description": "Manufacturing, production and services indices (measuring total economic output) and productivity (measuring efficiency, expressed as a ratio of output to input over a given period of time, for example output per person per hour).",
        "title": "Economic output and productivity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/economicoutputandproductivity",
                "id": "economicoutputandproductivity"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/economicoutputandproductivity/subtopics"
            }
        },
        "subtopics_ids": [
            "output",
            "productivitymeasures",
            "publicservicesproductivity"
        ]
    }
})
db.topics.insertOne({
    "id": "output",
    "next": {
        "id": "output",
        "description": "Economic output and activity of the UK. Includes manufacturing, production and services, and other measures of economic activity.",
        "title": "Output",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/output",
                "id": "output"
            },
            "content": {
                "href": "http://localhost:25300/topics/output/content"
            }
        }
    },
    "current": {
        "id": "output",
        "description": "Economic output and activity of the UK. Includes manufacturing, production and services, and other measures of economic activity.",
        "title": "Output",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/output",
                "id": "output"
            },
            "content": {
                "href": "http://localhost:25300/topics/output/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "productivitymeasures",
    "next": {
        "id": "productivitymeasures",
        "description": "Economic productivity measures, including output per hour, output per job and output per worker for the whole economy and a range of industries; productivity in the public sector; and international comparisons of productivity across the G7 nations.",
        "title": "Productivity measures",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/productivitymeasures",
                "id": "productivitymeasures"
            },
            "content": {
                "href": "http://localhost:25300/topics/productivitymeasures/content"
            }
        }
    },
    "current": {
        "id": "productivitymeasures",
        "description": "Economic productivity measures, including output per hour, output per job and output per worker for the whole economy and a range of industries; productivity in the public sector; and international comparisons of productivity across the G7 nations.",
        "title": "Productivity measures",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/productivitymeasures",
                "id": "productivitymeasures"
            },
            "content": {
                "href": "http://localhost:25300/topics/productivitymeasures/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "publicservicesproductivity",
    "next": {
        "id": "publicservicesproductivity",
        "description": "** no description **",
        "title": "Public services productivity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicservicesproductivity",
                "id": "publicservicesproductivity"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicservicesproductivity/content"
            }
        }
    },
    "current": {
        "id": "publicservicesproductivity",
        "description": "** no description **",
        "title": "Public services productivity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicservicesproductivity",
                "id": "publicservicesproductivity"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicservicesproductivity/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "environmentalaccounts",
    "next": {
        "id": "environmentalaccounts",
        "description": "Environmental accounts show how the environment contributes to the economy (for example, through the extraction of raw materials), the impacts that the economy has on the environment (for example, energy consumption and air emissions), and how society responds to environmental issues (for example, through taxation and expenditure on environmental protection).",
        "title": "Environmental accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/environmentalaccounts",
                "id": "environmentalaccounts"
            },
            "content": {
                "href": "http://localhost:25300/topics/environmentalaccounts/content"
            }
        }
    },
    "current": {
        "id": "environmentalaccounts",
        "description": "Environmental accounts show how the environment contributes to the economy (for example, through the extraction of raw materials), the impacts that the economy has on the environment (for example, energy consumption and air emissions), and how society responds to environmental issues (for example, through taxation and expenditure on environmental protection).",
        "title": "Environmental accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/environmentalaccounts",
                "id": "environmentalaccounts"
            },
            "content": {
                "href": "http://localhost:25300/topics/environmentalaccounts/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "governmentpublicsectorandtaxes",
    "next": {
        "id": "governmentpublicsectorandtaxes",
        "description": "Public sector spending, tax revenues and investments for the UK, including government debt and deficit (the gap between revenue and spending), research and development, and the effect of taxes.",
        "title": "Government, public sector and taxes",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/governmentpublicsectorandtaxes",
                "id": "governmentpublicsectorandtaxes"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/governmentpublicsectorandtaxes/subtopics"
            }
        },
        "subtopics_ids": [
            "localgovernmentfinance",
            "publicsectorfinance",
            "publicspending",
            "researchanddevelopmentexpenditure",
            "taxesandrevenue"
        ]
    },
    "current": {
        "id": "governmentpublicsectorandtaxes",
        "description": "Public sector spending, tax revenues and investments for the UK, including government debt and deficit (the gap between revenue and spending), research and development, and the effect of taxes.",
        "title": "Government, public sector and taxes",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/governmentpublicsectorandtaxes",
                "id": "governmentpublicsectorandtaxes"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/governmentpublicsectorandtaxes/subtopics"
            }
        },
        "subtopics_ids": [
            "localgovernmentfinance",
            "publicsectorfinance",
            "publicspending",
            "researchanddevelopmentexpenditure",
            "taxesandrevenue"
        ]
    }
})
db.topics.insertOne({
    "id": "localgovernmentfinance",
    "next": {
        "id": "localgovernmentfinance",
        "description": "This is what local government is given from the overall budget, and what they spend it on.",
        "title": "Local government finance",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/localgovernmentfinance",
                "id": "localgovernmentfinance"
            },
            "content": {
                "href": "http://localhost:25300/topics/localgovernmentfinance/content"
            }
        }
    },
    "current": {
        "id": "localgovernmentfinance",
        "description": "This is what local government is given from the overall budget, and what they spend it on.",
        "title": "Local government finance",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/localgovernmentfinance",
                "id": "localgovernmentfinance"
            },
            "content": {
                "href": "http://localhost:25300/topics/localgovernmentfinance/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "publicsectorfinance",
    "next": {
        "id": "publicsectorfinance",
        "description": "UK public sector spending, tax revenues and investments, including government debt and deficit (the gap between revenue and spending) and data submitted to the European Commission.",
        "title": "Public sector finance",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicsectorfinance",
                "id": "publicsectorfinance"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicsectorfinance/content"
            }
        }
    },
    "current": {
        "id": "publicsectorfinance",
        "description": "UK public sector spending, tax revenues and investments, including government debt and deficit (the gap between revenue and spending) and data submitted to the European Commission.",
        "title": "Public sector finance",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicsectorfinance",
                "id": "publicsectorfinance"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicsectorfinance/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "publicspending",
    "next": {
        "id": "publicspending",
        "description": "We look at what the UK government spends on the general public, such as pensions, provisions, and infrastructure.",
        "title": "Public spending",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicspending",
                "id": "publicspending"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicspending/content"
            }
        }
    },
    "current": {
        "id": "publicspending",
        "description": "We look at what the UK government spends on the general public, such as pensions, provisions, and infrastructure.",
        "title": "Public spending",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicspending",
                "id": "publicspending"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicspending/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "researchanddevelopmentexpenditure",
    "next": {
        "id": "researchanddevelopmentexpenditure",
        "description": "Research and development in the UK carried out or funded by business enterprises, higher education, government (including research councils) and private non-profit organisations.",
        "title": "Research and development expenditure",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/researchanddevelopmentexpenditure",
                "id": "researchanddevelopmentexpenditure"
            },
            "content": {
                "href": "http://localhost:25300/topics/researchanddevelopmentexpenditure/content"
            }
        }
    },
    "current": {
        "id": "researchanddevelopmentexpenditure",
        "description": "Research and development in the UK carried out or funded by business enterprises, higher education, government (including research councils) and private non-profit organisations.",
        "title": "Research and development expenditure",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/researchanddevelopmentexpenditure",
                "id": "researchanddevelopmentexpenditure"
            },
            "content": {
                "href": "http://localhost:25300/topics/researchanddevelopmentexpenditure/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "taxesandrevenue",
    "next": {
        "id": "taxesandrevenue",
        "description": "** no description **",
        "title": "Taxes and revenue",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/taxesandrevenue",
                "id": "taxesandrevenue"
            },
            "content": {
                "href": "http://localhost:25300/topics/taxesandrevenue/content"
            }
        }
    },
    "current": {
        "id": "taxesandrevenue",
        "description": "** no description **",
        "title": "Taxes and revenue",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/taxesandrevenue",
                "id": "taxesandrevenue"
            },
            "content": {
                "href": "http://localhost:25300/topics/taxesandrevenue/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "grossdomesticproductgdp",
    "next": {
        "id": "grossdomesticproductgdp",
        "description": "Estimates of GDP are released on a monthly and quarterly basis. Monthly estimates are released alongside other short-term economic indicators. The two quarterly estimates contain data from all three approaches to measuring GDP and are called the First quarterly estimate of GDP and the Quarterly National Accounts. Data sources feeding into the two types of releases are consistent with each other.",
        "title": "Gross Domestic Product (GDP)",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/grossdomesticproductgdp",
                "id": "grossdomesticproductgdp"
            },
            "content": {
                "href": "http://localhost:25300/topics/grossdomesticproductgdp/content"
            }
        }
    },
    "current": {
        "id": "grossdomesticproductgdp",
        "description": "Estimates of GDP are released on a monthly and quarterly basis. Monthly estimates are released alongside other short-term economic indicators. The two quarterly estimates contain data from all three approaches to measuring GDP and are called the First quarterly estimate of GDP and the Quarterly National Accounts. Data sources feeding into the two types of releases are consistent with each other.",
        "title": "Gross Domestic Product (GDP)",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/grossdomesticproductgdp",
                "id": "grossdomesticproductgdp"
            },
            "content": {
                "href": "http://localhost:25300/topics/grossdomesticproductgdp/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "grossvalueaddedgva",
    "next": {
        "id": "grossvalueaddedgva",
        "description": "Regional gross value added using production (GVA(P)) and income (GVA(I)) approaches. Regional gross value added is the value generated by any unit engaged in the production of goods and services. GVA per head is a useful way of comparing regions of different sizes. It is not, however, a measure of regional productivity.",
        "title": "Gross Value Added (GVA)",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/grossvalueaddedgva",
                "id": "grossvalueaddedgva"
            },
            "content": {
                "href": "http://localhost:25300/topics/grossvalueaddedgva/content"
            }
        }
    },
    "current": {
        "id": "grossvalueaddedgva",
        "description": "Regional gross value added using production (GVA(P)) and income (GVA(I)) approaches. Regional gross value added is the value generated by any unit engaged in the production of goods and services. GVA per head is a useful way of comparing regions of different sizes. It is not, however, a measure of regional productivity.",
        "title": "Gross Value Added (GVA)",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/grossvalueaddedgva",
                "id": "grossvalueaddedgva"
            },
            "content": {
                "href": "http://localhost:25300/topics/grossvalueaddedgva/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "inflationandpriceindices",
    "next": {
        "id": "inflationandpriceindices",
        "description": "The rate of increase in prices for goods and services. Measures of inflation and prices include consumer price inflation, producer price inflation, the house price index, index of private housing rental prices, and construction output price indices. ",
        "title": "Inflation and price indices",
        "keywords": [
            "Consumer Price Index,Retail Price Index,Producer Price Index,Services Producer Price Indices,Index of Private Housing Rental Prices,CPIH,RPIJ"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/inflationandpriceindices",
                "id": "inflationandpriceindices"
            },
            "content": {
                "href": "http://localhost:25300/topics/inflationandpriceindices/content"
            }
        }
    },
    "current": {
        "id": "inflationandpriceindices",
        "description": "The rate of increase in prices for goods and services. Measures of inflation and prices include consumer price inflation, producer price inflation, the house price index, index of private housing rental prices, and construction output price indices. ",
        "title": "Inflation and price indices",
        "keywords": [
            "Consumer Price Index,Retail Price Index,Producer Price Index,Services Producer Price Indices,Index of Private Housing Rental Prices,CPIH,RPIJ"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/inflationandpriceindices",
                "id": "inflationandpriceindices"
            },
            "content": {
                "href": "http://localhost:25300/topics/inflationandpriceindices/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "investmentspensionsandtrusts",
    "next": {
        "id": "investmentspensionsandtrusts",
        "description": "Net flows of investment into the UK, the number of people who hold pensions of different types, and investments made by various types of trusts. ",
        "title": "Investments, pensions and trusts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/investmentspensionsandtrusts",
                "id": "investmentspensionsandtrusts"
            },
            "content": {
                "href": "http://localhost:25300/topics/investmentspensionsandtrusts/content"
            }
        }
    },
    "current": {
        "id": "investmentspensionsandtrusts",
        "description": "Net flows of investment into the UK, the number of people who hold pensions of different types, and investments made by various types of trusts. ",
        "title": "Investments, pensions and trusts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/investmentspensionsandtrusts",
                "id": "investmentspensionsandtrusts"
            },
            "content": {
                "href": "http://localhost:25300/topics/investmentspensionsandtrusts/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "nationalaccounts",
    "next": {
        "id": "nationalaccounts",
        "description": "Core accounts for the UK economy as a whole; individual sectors (sector accounts); accounts for the regions, subregions and local areas of the UK; and satellite accounts that cover activities linked to the economy. The national accounts framework brings units and transactions together to provide a simple and understandable description of production, income, consumption, accumulation and wealth.",
        "title": "National accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/nationalaccounts",
                "id": "nationalaccounts"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/nationalaccounts/subtopics"
            }
        },
        "subtopics_ids": [
            "balanceofpayments",
            "satelliteaccounts",
            "supplyandusetables",
            "uksectoraccounts"
        ],
        "spotlight": [
            {
                "href": "/economy/grossdomesticproductgdp/compendium/unitedkingdomnationalaccountsthebluebook/latest",
                "title": "UK National Accounts, The Blue Book"
            }
        ]
    },
    "current": {
        "id": "nationalaccounts",
        "description": "Core accounts for the UK economy as a whole; individual sectors (sector accounts); accounts for the regions, subregions and local areas of the UK; and satellite accounts that cover activities linked to the economy. The national accounts framework brings units and transactions together to provide a simple and understandable description of production, income, consumption, accumulation and wealth.",
        "title": "National accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/nationalaccounts",
                "id": "nationalaccounts"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/nationalaccounts/subtopics"
            }
        },
        "subtopics_ids": [
            "balanceofpayments",
            "satelliteaccounts",
            "supplyandusetables",
            "uksectoraccounts"
        ],
        "spotlight": [
            {
                "href": "/economy/grossdomesticproductgdp/compendium/unitedkingdomnationalaccountsthebluebook/latest",
                "title": "UK National Accounts, The Blue Book"
            }
        ]
    }
})
db.topics.insertOne({
    "id": "balanceofpayments",
    "next": {
        "id": "balanceofpayments",
        "description": "All economic transactions between residents of the UK and the rest of the world.",
        "title": "Balance of payments",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/balanceofpayments",
                "id": "balanceofpayments"
            },
            "content": {
                "href": "http://localhost:25300/topics/balanceofpayments/content"
            }
        }
    },
    "current": {
        "id": "balanceofpayments",
        "description": "All economic transactions between residents of the UK and the rest of the world.",
        "title": "Balance of payments",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/balanceofpayments",
                "id": "balanceofpayments"
            },
            "content": {
                "href": "http://localhost:25300/topics/balanceofpayments/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "satelliteaccounts",
    "next": {
        "id": "satelliteaccounts",
        "description": "Accounts that cover activities linked to the economy but not part of the core UK national accounts including environmental accounts, tourism satellite account, and household satellite accounts.",
        "title": "Satellite accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/satelliteaccounts",
                "id": "satelliteaccounts"
            },
            "content": {
                "href": "http://localhost:25300/topics/satelliteaccounts/content"
            }
        }
    },
    "current": {
        "id": "satelliteaccounts",
        "description": "Accounts that cover activities linked to the economy but not part of the core UK national accounts including environmental accounts, tourism satellite account, and household satellite accounts.",
        "title": "Satellite accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/satelliteaccounts",
                "id": "satelliteaccounts"
            },
            "content": {
                "href": "http://localhost:25300/topics/satelliteaccounts/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "supplyandusetables",
    "next": {
        "id": "supplyandusetables",
        "description": "Balances showing the relationship between components of value added, industry inputs and outputs, and product supply and demand. These tables are a source for the data underlying Gross Domestic Product. ",
        "title": "Supply and use tables",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/supplyandusetables",
                "id": "supplyandusetables"
            },
            "content": {
                "href": "http://localhost:25300/topics/supplyandusetables/content"
            }
        }
    },
    "current": {
        "id": "supplyandusetables",
        "description": "Balances showing the relationship between components of value added, industry inputs and outputs, and product supply and demand. These tables are a source for the data underlying Gross Domestic Product. ",
        "title": "Supply and use tables",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/supplyandusetables",
                "id": "supplyandusetables"
            },
            "content": {
                "href": "http://localhost:25300/topics/supplyandusetables/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "uksectoraccounts",
    "next": {
        "id": "uksectoraccounts",
        "description": "Transactions of particular groups of institutions (sectors) within the UK economy, showing how the income from production is distributed and redistributed and how savings are used to add wealth through investment in physical or financial assets.",
        "title": "UK sector accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/uksectoraccounts",
                "id": "uksectoraccounts"
            },
            "content": {
                "href": "http://localhost:25300/topics/uksectoraccounts/content"
            }
        }
    },
    "current": {
        "id": "uksectoraccounts",
        "description": "Transactions of particular groups of institutions (sectors) within the UK economy, showing how the income from production is distributed and redistributed and how savings are used to add wealth through investment in physical or financial assets.",
        "title": "UK sector accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/uksectoraccounts",
                "id": "uksectoraccounts"
            },
            "content": {
                "href": "http://localhost:25300/topics/uksectoraccounts/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "regionalaccounts",
    "next": {
        "id": "regionalaccounts",
        "description": "Accounts for regions, sub-regions and local areas of the UK. These accounts allow comparisons between regions and against a UK average. Statistics include regional gross value added (GVA) and figures on regional gross disposable household income (GDHI).",
        "title": "Regional accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/regionalaccounts",
                "id": "regionalaccounts"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/regionalaccounts/subtopics"
            }
        },
        "subtopics_ids": [
            "grossdisposablehouseholdincome"
        ]
    },
    "current": {
        "id": "regionalaccounts",
        "description": "Accounts for regions, sub-regions and local areas of the UK. These accounts allow comparisons between regions and against a UK average. Statistics include regional gross value added (GVA) and figures on regional gross disposable household income (GDHI).",
        "title": "Regional accounts",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/regionalaccounts",
                "id": "regionalaccounts"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/regionalaccounts/subtopics"
            }
        },
        "subtopics_ids": [
            "grossdisposablehouseholdincome"
        ]
    }
})
db.topics.insertOne({
    "id": "grossdisposablehouseholdincome",
    "next": {
        "id": "grossdisposablehouseholdincome",
        "description": "The amount of money that that all of the individuals in the household sector have available for spending or saving after income distribution measures (for example, taxes, social contributions and benefits) have taken effect. GDHI does not provide measures relating to actual households or family units. The figures cover regions, sub-regions and local areas of the UK. ",
        "title": "Gross disposable household income",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/grossdisposablehouseholdincome",
                "id": "grossdisposablehouseholdincome"
            },
            "content": {
                "href": "http://localhost:25300/topics/grossdisposablehouseholdincome/content"
            }
        }
    },
    "current": {
        "id": "grossdisposablehouseholdincome",
        "description": "The amount of money that that all of the individuals in the household sector have available for spending or saving after income distribution measures (for example, taxes, social contributions and benefits) have taken effect. GDHI does not provide measures relating to actual households or family units. The figures cover regions, sub-regions and local areas of the UK. ",
        "title": "Gross disposable household income",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/grossdisposablehouseholdincome",
                "id": "grossdisposablehouseholdincome"
            },
            "content": {
                "href": "http://localhost:25300/topics/grossdisposablehouseholdincome/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "employmentandlabourmarket",
    "next": {
        "id": "employmentandlabourmarket",
        "description": "People in and out of work covering employment, unemployment, types of work, earnings, working patterns and workplace disputes.",
        "title": "Employment and labour market",
        "keywords": [
            "economic activity",
            "jobs",
            "vacancies",
            "hours of work",
            "unemployment"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/employmentandlabourmarket",
                "id": "employmentandlabourmarket"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/employmentandlabourmarket/subtopics"
            }
        },
        "subtopics_ids": [
            "peopleinwork",
            "peoplenotinwork"
        ]
    },
    "current": {
        "id": "employmentandlabourmarket",
        "description": "People in and out of work covering employment, unemployment, types of work, earnings, working patterns and workplace disputes.",
        "title": "Employment and labour market",
        "keywords": [
            "economic activity",
            "jobs",
            "vacancies",
            "hours of work",
            "unemployment"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/employmentandlabourmarket",
                "id": "employmentandlabourmarket"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/employmentandlabourmarket/subtopics"
            }
        },
        "subtopics_ids": [
            "peopleinwork",
            "peoplenotinwork"
        ]
    }
})
db.topics.insertOne({
    "id": "peopleinwork",
    "next": {
        "id": "peopleinwork",
        "description": "Employment data covering employment rates, hours of work, claimants and earnings.",
        "title": "People in work",
        "keywords": [
            "gender pay gap",
            "average weekly earnings",
            "employment rates",
            "pensions",
            "labour productivity"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/peopleinwork",
                "id": "peopleinwork"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/peopleinwork/subtopics"
            }
        },
        "subtopics_ids": [
            "earningsandworkinghours",
            "employmentandemployeetypes",
            "labourproductivity",
            "publicsectorpersonnel",
            "workplacedisputesandworkingconditions",
            "workplacepensions"
        ],
        "spotlight": [
            {
                "href": "/employmentandlabourmarket/peopleinwork/employmentandemployeetypes/bulletins/uklabourmarket/latest",
                "title": "Labour market overview, UK"
            },
            {
                "href": "/employmentandlabourmarket/peopleinwork/earningsandworkinghours/bulletins/annualsurveyofhoursandearnings/latest",
                "title": "Employee earnings in the UK"
            }
        ]
    },
    "current": {
        "id": "peopleinwork",
        "description": "Employment data covering employment rates, hours of work, claimants and earnings.",
        "title": "People in work",
        "keywords": [
            "gender pay gap",
            "average weekly earnings",
            "employment rates",
            "pensions",
            "labour productivity"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/peopleinwork",
                "id": "peopleinwork"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/peopleinwork/subtopics"
            }
        },
        "subtopics_ids": [
            "earningsandworkinghours",
            "employmentandemployeetypes",
            "labourproductivity",
            "publicsectorpersonnel",
            "workplacedisputesandworkingconditions",
            "workplacepensions"
        ],
        "spotlight": [
            {
                "href": "/employmentandlabourmarket/peopleinwork/employmentandemployeetypes/bulletins/uklabourmarket/latest",
                "title": "Labour market overview, UK"
            },
            {
                "href": "/employmentandlabourmarket/peopleinwork/earningsandworkinghours/bulletins/annualsurveyofhoursandearnings/latest",
                "title": "Employee earnings in the UK"
            }
        ]
    }
})
db.topics.insertOne({
    "id": "earningsandworkinghours",
    "next": {
        "id": "earningsandworkinghours",
        "description": "Average weekly earnings of people in the UK and information on the gender pay gap and low pay. Data from Average Weekly Earnings (AWE) and the Annual Survey of Hours and Earnings (ASHE).",
        "title": "Earnings and working hours",
        "keywords": [
            "employee earnings",
            "labour market",
            "hours worked",
            "minimum wage",
            "average salary"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/earningsandworkinghours",
                "id": "earningsandworkinghours"
            },
            "content": {
                "href": "http://localhost:25300/topics/earningsandworkinghours/content"
            }
        }
    },
    "current": {
        "id": "earningsandworkinghours",
        "description": "Average weekly earnings of people in the UK and information on the gender pay gap and low pay. Data from Average Weekly Earnings (AWE) and the Annual Survey of Hours and Earnings (ASHE).",
        "title": "Earnings and working hours",
        "keywords": [
            "employee earnings",
            "labour market",
            "hours worked",
            "minimum wage",
            "average salary"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/earningsandworkinghours",
                "id": "earningsandworkinghours"
            },
            "content": {
                "href": "http://localhost:25300/topics/earningsandworkinghours/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "employmentandemployeetypes",
    "next": {
        "id": "employmentandemployeetypes",
        "description": "Employment rates show the number of people in paid work as a proportion of the population, broken down by age and sex. Includes information on the number of people in employment and vacancies.",
        "title": "Employment and employee types",
        "keywords": [
            "jobs",
            "labour market",
            "occupation",
            "full-time employment",
            "part-time employment"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/employmentandemployeetypes",
                "id": "employmentandemployeetypes"
            },
            "content": {
                "href": "http://localhost:25300/topics/employmentandemployeetypes/content"
            }
        }
    },
    "current": {
        "id": "employmentandemployeetypes",
        "description": "Employment rates show the number of people in paid work as a proportion of the population, broken down by age and sex. Includes information on the number of people in employment and vacancies.",
        "title": "Employment and employee types",
        "keywords": [
            "jobs",
            "labour market",
            "occupation",
            "full-time employment",
            "part-time employment"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/employmentandemployeetypes",
                "id": "employmentandemployeetypes"
            },
            "content": {
                "href": "http://localhost:25300/topics/employmentandemployeetypes/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "labourproductivity",
    "next": {
        "id": "labourproductivity",
        "description": "Efficiency of the UK workforce, including output per worker, per job and per hour. Data are available by industry and by region.",
        "title": "Labour productivity",
        "keywords": [
            "economic growth",
            "real GDP",
            "UK productivity",
            "productivity puzzle"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/labourproductivity",
                "id": "labourproductivity"
            },
            "content": {
                "href": "http://localhost:25300/topics/labourproductivity/content"
            }
        }
    },
    "current": {
        "id": "labourproductivity",
        "description": "Efficiency of the UK workforce, including output per worker, per job and per hour. Data are available by industry and by region.",
        "title": "Labour productivity",
        "keywords": [
            "economic growth",
            "real GDP",
            "UK productivity",
            "productivity puzzle"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/labourproductivity",
                "id": "labourproductivity"
            },
            "content": {
                "href": "http://localhost:25300/topics/labourproductivity/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "publicsectorpersonnel",
    "next": {
        "id": "publicsectorpersonnel",
        "description": "People employed in central and local government, and public corporations, including second jobs in the public sector. Includes Civil Service employment with regional and diversity analyses.",
        "title": "Public sector personnel",
        "keywords": [
            "public sector employment",
            "civil servants",
            "public sector workers",
            "labour market",
            "central government"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicsectorpersonnel",
                "id": "publicsectorpersonnel"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicsectorpersonnel/content"
            }
        }
    },
    "current": {
        "id": "publicsectorpersonnel",
        "description": "People employed in central and local government, and public corporations, including second jobs in the public sector. Includes Civil Service employment with regional and diversity analyses.",
        "title": "Public sector personnel",
        "keywords": [
            "public sector employment",
            "civil servants",
            "public sector workers",
            "labour market",
            "central government"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/publicsectorpersonnel",
                "id": "publicsectorpersonnel"
            },
            "content": {
                "href": "http://localhost:25300/topics/publicsectorpersonnel/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "workplacedisputesandworkingconditions",
    "next": {
        "id": "workplacedisputesandworkingconditions",
        "description": "Work stoppages because of disputes between employers and employees. Includes strikes and lock-outs, number of days lost in the public and private sectors, and number of workers involved.",
        "title": "Workplace disputes and working conditions",
        "keywords": [
            "working days lost",
            "strike action",
            "industrial action",
            "trade union",
            "labour disputes"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/workplacedisputesandworkingconditions",
                "id": "workplacedisputesandworkingconditions"
            },
            "content": {
                "href": "http://localhost:25300/topics/workplacedisputesandworkingconditions/content"
            }
        }
    },
    "current": {
        "id": "workplacedisputesandworkingconditions",
        "description": "Work stoppages because of disputes between employers and employees. Includes strikes and lock-outs, number of days lost in the public and private sectors, and number of workers involved.",
        "title": "Workplace disputes and working conditions",
        "keywords": [
            "working days lost",
            "strike action",
            "industrial action",
            "trade union",
            "labour disputes"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/workplacedisputesandworkingconditions",
                "id": "workplacedisputesandworkingconditions"
            },
            "content": {
                "href": "http://localhost:25300/topics/workplacedisputesandworkingconditions/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "workplacepensions",
    "next": {
        "id": "workplacepensions",
        "description": "Pensions linked to a person's workplace, including defined ambition (DA), defined benefit (DB) and defined contribution (DC) schemes.",
        "title": "Workplace pensions",
        "keywords": [
            "pension contributions",
            "state pension",
            "employer contributions",
            "public sector",
            "pension membership"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/workplacepensions",
                "id": "workplacepensions"
            },
            "content": {
                "href": "http://localhost:25300/topics/workplacepensions/content"
            }
        }
    },
    "current": {
        "id": "workplacepensions",
        "description": "Pensions linked to a person's workplace, including defined ambition (DA), defined benefit (DB) and defined contribution (DC) schemes.",
        "title": "Workplace pensions",
        "keywords": [
            "pension contributions",
            "state pension",
            "employer contributions",
            "public sector",
            "pension membership"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/workplacepensions",
                "id": "workplacepensions"
            },
            "content": {
                "href": "http://localhost:25300/topics/workplacepensions/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "peoplenotinwork",
    "next": {
        "id": "peoplenotinwork",
        "description": "Unemployed and economically inactive people in the UK including claimants of out-of-work benefits and the number of redundancies.\n",
        "title": "People not in work",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/peoplenotinwork",
                "id": "peoplenotinwork"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/peoplenotinwork/subtopics"
            }
        },
        "subtopics_ids": [
            "economicinactivity",
            "outofworkbenefits",
            "redundancies",
            "unemployment"
        ],
        "spotlight": [
            {
                "href": "/employmentandlabourmarket/peopleinwork/employmentandemployeetypes/bulletins/uklabourmarket/latest",
                "title": "Labour market overview, UK"
            }
        ]
    },
    "current": {
        "id": "peoplenotinwork",
        "description": "Unemployed and economically inactive people in the UK including claimants of out-of-work benefits and the number of redundancies.\n",
        "title": "People not in work",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/peoplenotinwork",
                "id": "peoplenotinwork"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/peoplenotinwork/subtopics"
            }
        },
        "subtopics_ids": [
            "economicinactivity",
            "outofworkbenefits",
            "redundancies",
            "unemployment"
        ],
        "spotlight": [
            {
                "href": "/employmentandlabourmarket/peopleinwork/employmentandemployeetypes/bulletins/uklabourmarket/latest",
                "title": "Labour market overview, UK"
            }
        ]
    }
})
db.topics.insertOne({
    "id": "economicinactivity",
    "next": {
        "id": "economicinactivity",
        "description": "People not in employment who have not been seeking work within the last 4 weeks and/or are unable to start work within the next 2 weeks.",
        "title": "Economic inactivity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/economicinactivity",
                "id": "economicinactivity"
            },
            "content": {
                "href": "http://localhost:25300/topics/economicinactivity/content"
            }
        }
    },
    "current": {
        "id": "economicinactivity",
        "description": "People not in employment who have not been seeking work within the last 4 weeks and/or are unable to start work within the next 2 weeks.",
        "title": "Economic inactivity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/economicinactivity",
                "id": "economicinactivity"
            },
            "content": {
                "href": "http://localhost:25300/topics/economicinactivity/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "outofworkbenefits",
    "next": {
        "id": "outofworkbenefits",
        "description": "Claimants of unemployment related benefits including Employment and Support Allowance and other incapacity benefits, and Income Support and Pension Credit. ",
        "title": "Out of work benefits",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/outofworkbenefits",
                "id": "outofworkbenefits"
            },
            "content": {
                "href": "http://localhost:25300/topics/outofworkbenefits/content"
            }
        }
    },
    "current": {
        "id": "outofworkbenefits",
        "description": "Claimants of unemployment related benefits including Employment and Support Allowance and other incapacity benefits, and Income Support and Pension Credit. ",
        "title": "Out of work benefits",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/outofworkbenefits",
                "id": "outofworkbenefits"
            },
            "content": {
                "href": "http://localhost:25300/topics/outofworkbenefits/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "redundancies",
    "next": {
        "id": "redundancies",
        "description": "People who have been made redundant or have taken voluntary redundancy.",
        "title": "Redundancies",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/redundancies",
                "id": "redundancies"
            },
            "content": {
                "href": "http://localhost:25300/topics/redundancies/content"
            }
        }
    },
    "current": {
        "id": "redundancies",
        "description": "People who have been made redundant or have taken voluntary redundancy.",
        "title": "Redundancies",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/redundancies",
                "id": "redundancies"
            },
            "content": {
                "href": "http://localhost:25300/topics/redundancies/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "unemployment",
    "next": {
        "id": "unemployment",
        "description": "The level and rate of UK unemployment measured by the Labour Force Survey (LFS), using the International Labour Organisation's definition of unemployment.",
        "title": "Unemployment",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/unemployment",
                "id": "unemployment"
            },
            "content": {
                "href": "http://localhost:25300/topics/unemployment/content"
            }
        }
    },
    "current": {
        "id": "unemployment",
        "description": "The level and rate of UK unemployment measured by the Labour Force Survey (LFS), using the International Labour Organisation's definition of unemployment.",
        "title": "Unemployment",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/unemployment",
                "id": "unemployment"
            },
            "content": {
                "href": "http://localhost:25300/topics/unemployment/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "peoplepopulationandcommunity",
    "next": {
        "id": "peoplepopulationandcommunity",
        "description": "People living in the UK, changes in the population, how we spend our money, and data on crime, relationships, health and religion. These statistics help us build a detailed picture of how we live.",
        "title": "People, population and community",
        "keywords": [
            "Crime",
            "Education",
            "Marriage",
            "Religion",
            "Immigration"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/peoplepopulationandcommunity",
                "id": "peoplepopulationandcommunity"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/peoplepopulationandcommunity/subtopics"
            }
        },
        "subtopics_ids": [
            "birthsdeathsandmarriages",
            "community",
            "crimeandjustice",
            "culturalidentity",
            "educationandchildcare",
            "elections",
            "healthandsocialcare",
            "householdcharacteristics",
            "housing",
            "leisureandtourism",
            "personalandhouseholdfinances",
            "populationandmigration",
            "wellbeing"
        ]
    },
    "current": {
        "id": "peoplepopulationandcommunity",
        "description": "People living in the UK, changes in the population, how we spend our money, and data on crime, relationships, health and religion. These statistics help us build a detailed picture of how we live.",
        "title": "People, population and community",
        "keywords": [
            "Crime",
            "Education",
            "Marriage",
            "Religion",
            "Immigration"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/peoplepopulationandcommunity",
                "id": "peoplepopulationandcommunity"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/peoplepopulationandcommunity/subtopics"
            }
        },
        "subtopics_ids": [
            "birthsdeathsandmarriages",
            "community",
            "crimeandjustice",
            "culturalidentity",
            "educationandchildcare",
            "elections",
            "healthandsocialcare",
            "householdcharacteristics",
            "housing",
            "leisureandtourism",
            "personalandhouseholdfinances",
            "populationandmigration",
            "wellbeing"
        ]
    }
})
db.topics.insertOne({
    "id": "birthsdeathsandmarriages",
    "next": {
        "id": "birthsdeathsandmarriages",
        "description": "Life events in the UK including fertility rates, live births and stillbirths, family composition, life expectancy and deaths. This tells us about the health and relationships of the population.",
        "title": "Births, deaths and marriages",
        "keywords": [
            "Mortality",
            "Families",
            "Life expectancies",
            "Babies",
            "Divorce"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/birthsdeathsandmarriages",
                "id": "birthsdeathsandmarriages"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/birthsdeathsandmarriages/subtopics"
            }
        },
        "subtopics_ids": [
            "adoption",
            "ageing",
            "conceptionandfertilityrates",
            "deaths",
            "divorce",
            "families",
            "lifeexpectancies",
            "livebirths",
            "marriagecohabitationandcivilpartnerships",
            "maternities",
            "stillbirths"
        ]
    },
    "current": {
        "id": "birthsdeathsandmarriages",
        "description": "Life events in the UK including fertility rates, live births and stillbirths, family composition, life expectancy and deaths. This tells us about the health and relationships of the population.",
        "title": "Births, deaths and marriages",
        "keywords": [
            "Mortality",
            "Families",
            "Life expectancies",
            "Babies",
            "Divorce"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/birthsdeathsandmarriages",
                "id": "birthsdeathsandmarriages"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/birthsdeathsandmarriages/subtopics"
            }
        },
        "subtopics_ids": [
            "adoption",
            "ageing",
            "conceptionandfertilityrates",
            "deaths",
            "divorce",
            "families",
            "lifeexpectancies",
            "livebirths",
            "marriagecohabitationandcivilpartnerships",
            "maternities",
            "stillbirths"
        ]
    }
})
db.topics.insertOne({
    "id": "adoption",
    "next": {
        "id": "adoption",
        "description": "Adoption registrations including the age and gender of the children being adopted and whether they were born inside marriage.",
        "title": "Adoption",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/adoption",
                "id": "adoption"
            },
            "content": {
                "href": "http://localhost:25300/topics/adoption/content"
            }
        }
    },
    "current": {
        "id": "adoption",
        "description": "Adoption registrations including the age and gender of the children being adopted and whether they were born inside marriage.",
        "title": "Adoption",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/adoption",
                "id": "adoption"
            },
            "content": {
                "href": "http://localhost:25300/topics/adoption/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "ageing",
    "next": {
        "id": "ageing",
        "description": "Estimates of those aged 90 years and over in the UK, including what determines life expectancy and indicators of the health of the very old. These statistics help us understand the needs of an ageing population. ",
        "title": "Ageing",
        "keywords": [
            "Life expectancy",
            "Population",
            "Elderly",
            "Retirement",
            "Pensioner"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/ageing",
                "id": "ageing"
            },
            "content": {
                "href": "http://localhost:25300/topics/ageing/content"
            }
        }
    },
    "current": {
        "id": "ageing",
        "description": "Estimates of those aged 90 years and over in the UK, including what determines life expectancy and indicators of the health of the very old. These statistics help us understand the needs of an ageing population. ",
        "title": "Ageing",
        "keywords": [
            "Life expectancy",
            "Population",
            "Elderly",
            "Retirement",
            "Pensioner"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/ageing",
                "id": "ageing"
            },
            "content": {
                "href": "http://localhost:25300/topics/ageing/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "conceptionandfertilityrates",
    "next": {
        "id": "conceptionandfertilityrates",
        "description": "Childbearing, conceptions, births and abortions in the UK. These statistics inform us about average family size and maternity trends among women, by age and area. ",
        "title": "Conception and fertility rates",
        "keywords": [
            "Births",
            "abortions",
            "babies",
            "children",
            "pregnancy"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/conceptionandfertilityrates",
                "id": "conceptionandfertilityrates"
            },
            "content": {
                "href": "http://localhost:25300/topics/conceptionandfertilityrates/content"
            }
        }
    },
    "current": {
        "id": "conceptionandfertilityrates",
        "description": "Childbearing, conceptions, births and abortions in the UK. These statistics inform us about average family size and maternity trends among women, by age and area. ",
        "title": "Conception and fertility rates",
        "keywords": [
            "Births",
            "abortions",
            "babies",
            "children",
            "pregnancy"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/conceptionandfertilityrates",
                "id": "conceptionandfertilityrates"
            },
            "content": {
                "href": "http://localhost:25300/topics/conceptionandfertilityrates/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "deaths",
    "next": {
        "id": "deaths",
        "description": "Deaths broken down by age, sex, area and cause of death. ",
        "title": "Deaths",
        "keywords": [
            "Mortality",
            "Suicide",
            "MRSA",
            "Depression",
            "drugs"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/deaths",
                "id": "deaths"
            },
            "content": {
                "href": "http://localhost:25300/topics/deaths/content"
            }
        }
    },
    "current": {
        "id": "deaths",
        "description": "Deaths broken down by age, sex, area and cause of death. ",
        "title": "Deaths",
        "keywords": [
            "Mortality",
            "Suicide",
            "MRSA",
            "Depression",
            "drugs"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/deaths",
                "id": "deaths"
            },
            "content": {
                "href": "http://localhost:25300/topics/deaths/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "divorce",
    "next": {
        "id": "divorce",
        "description": "Divorces taking place covering dissolutions and annulments of marriage by previous marital status, sex and age of persons divorcing, children of  divorced couples, fact proven at divorce and to whom granted.",
        "title": "Divorce",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/divorce",
                "id": "divorce"
            },
            "content": {
                "href": "http://localhost:25300/topics/divorce/content"
            }
        }
    },
    "current": {
        "id": "divorce",
        "description": "Divorces taking place covering dissolutions and annulments of marriage by previous marital status, sex and age of persons divorcing, children of  divorced couples, fact proven at divorce and to whom granted.",
        "title": "Divorce",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/divorce",
                "id": "divorce"
            },
            "content": {
                "href": "http://localhost:25300/topics/divorce/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "families",
    "next": {
        "id": "families",
        "description": "The composition of families and households, including data on lone parents, married couples and civil partnership families. Household size and household types, including people living alone, multi-family households and households where members are all unrelated are also provided.",
        "title": "Families",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/families",
                "id": "families"
            },
            "content": {
                "href": "http://localhost:25300/topics/families/content"
            }
        }
    },
    "current": {
        "id": "families",
        "description": "The composition of families and households, including data on lone parents, married couples and civil partnership families. Household size and household types, including people living alone, multi-family households and households where members are all unrelated are also provided.",
        "title": "Families",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/families",
                "id": "families"
            },
            "content": {
                "href": "http://localhost:25300/topics/families/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "lifeexpectancies",
    "next": {
        "id": "lifeexpectancies",
        "description": "How long, on average, people can expect to live using estimates of the population and the number of deaths.",
        "title": "Life expectancies",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/lifeexpectancies",
                "id": "lifeexpectancies"
            },
            "content": {
                "href": "http://localhost:25300/topics/lifeexpectancies/content"
            }
        }
    },
    "current": {
        "id": "lifeexpectancies",
        "description": "How long, on average, people can expect to live using estimates of the population and the number of deaths.",
        "title": "Life expectancies",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/lifeexpectancies",
                "id": "lifeexpectancies"
            },
            "content": {
                "href": "http://localhost:25300/topics/lifeexpectancies/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "livebirths",
    "next": {
        "id": "livebirths",
        "description": "Live births by age of mother/father, sex, marital status, country of birth, socio-economic status, previous children and area.",
        "title": "Live births",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/livebirths",
                "id": "livebirths"
            },
            "content": {
                "href": "http://localhost:25300/topics/livebirths/content"
            }
        }
    },
    "current": {
        "id": "livebirths",
        "description": "Live births by age of mother/father, sex, marital status, country of birth, socio-economic status, previous children and area.",
        "title": "Live births",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/livebirths",
                "id": "livebirths"
            },
            "content": {
                "href": "http://localhost:25300/topics/livebirths/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "marriagecohabitationandcivilpartnerships",
    "next": {
        "id": "marriagecohabitationandcivilpartnerships",
        "description": "Marriages formed, civil partnerships formed and dissolved, and estimates of the population by marital status and living arrangements.",
        "title": "Marriage, cohabitation and civil partnerships",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/marriagecohabitationandcivilpartnerships",
                "id": "marriagecohabitationandcivilpartnerships"
            },
            "content": {
                "href": "http://localhost:25300/topics/marriagecohabitationandcivilpartnerships/content"
            }
        }
    },
    "current": {
        "id": "marriagecohabitationandcivilpartnerships",
        "description": "Marriages formed, civil partnerships formed and dissolved, and estimates of the population by marital status and living arrangements.",
        "title": "Marriage, cohabitation and civil partnerships",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/marriagecohabitationandcivilpartnerships",
                "id": "marriagecohabitationandcivilpartnerships"
            },
            "content": {
                "href": "http://localhost:25300/topics/marriagecohabitationandcivilpartnerships/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "maternities",
    "next": {
        "id": "maternities",
        "description": "Women having babies (including stillbirths). A maternity is a pregnancy resulting in the birth of 1 or more children, therefore, these figures are not the same as the number of babies born.",
        "title": "Maternities",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/maternities",
                "id": "maternities"
            },
            "content": {
                "href": "http://localhost:25300/topics/maternities/content"
            }
        }
    },
    "current": {
        "id": "maternities",
        "description": "Women having babies (including stillbirths). A maternity is a pregnancy resulting in the birth of 1 or more children, therefore, these figures are not the same as the number of babies born.",
        "title": "Maternities",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/maternities",
                "id": "maternities"
            },
            "content": {
                "href": "http://localhost:25300/topics/maternities/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "stillbirths",
    "next": {
        "id": "stillbirths",
        "description": "Stillbirths containing data on cause of death and sex, plus analyses by some of the key risk factors affecting stillbirths, including age of mother and birthweight.",
        "title": "Stillbirths",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/stillbirths",
                "id": "stillbirths"
            },
            "content": {
                "href": "http://localhost:25300/topics/stillbirths/content"
            }
        }
    },
    "current": {
        "id": "stillbirths",
        "description": "Stillbirths containing data on cause of death and sex, plus analyses by some of the key risk factors affecting stillbirths, including age of mother and birthweight.",
        "title": "Stillbirths",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/stillbirths",
                "id": "stillbirths"
            },
            "content": {
                "href": "http://localhost:25300/topics/stillbirths/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "crimeandjustice",
    "next": {
        "id": "crimeandjustice",
        "description": "Crimes committed and the victims' characteristics, sourced from crimes recorded by the police and from the Crime Survey for England and Wales (CSEW). The outcomes of crime in terms of what happened to the offender are also included.",
        "title": "Crime and justice",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/crimeandjustice",
                "id": "crimeandjustice"
            },
            "content": {
                "href": "http://localhost:25300/topics/crimeandjustice/content"
            }
        }
    },
    "current": {
        "id": "crimeandjustice",
        "description": "Crimes committed and the victims' characteristics, sourced from crimes recorded by the police and from the Crime Survey for England and Wales (CSEW). The outcomes of crime in terms of what happened to the offender are also included.",
        "title": "Crime and justice",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/crimeandjustice",
                "id": "crimeandjustice"
            },
            "content": {
                "href": "http://localhost:25300/topics/crimeandjustice/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "culturalidentity",
    "next": {
        "id": "culturalidentity",
        "description": "How people in the UK see themselves today in terms of ethnicity, sexual identity, religion and language, and how this has changed over time. We use a diverse range of sources for this data.",
        "title": "Cultural identity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/culturalidentity",
                "id": "culturalidentity"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/culturalidentity/subtopics"
            }
        },
        "subtopics_ids": [
            "ethnicity",
            "language",
            "religion",
            "sexuality"
        ]
    },
    "current": {
        "id": "culturalidentity",
        "description": "How people in the UK see themselves today in terms of ethnicity, sexual identity, religion and language, and how this has changed over time. We use a diverse range of sources for this data.",
        "title": "Cultural identity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/culturalidentity",
                "id": "culturalidentity"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/culturalidentity/subtopics"
            }
        },
        "subtopics_ids": [
            "ethnicity",
            "language",
            "religion",
            "sexuality"
        ]
    }
})
db.topics.insertOne({
    "id": "ethnicity",
    "next": {
        "id": "ethnicity",
        "description": "Analyses include ethnic identities among the non-UK born population in England and Wales, labour market participation, trends in general health and unpaid care by ethnic group, and inter-ethnic relationships. ",
        "title": "Ethnicity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/ethnicity",
                "id": "ethnicity"
            },
            "content": {
                "href": "http://localhost:25300/topics/ethnicity/content"
            }
        }
    },
    "current": {
        "id": "ethnicity",
        "description": "Analyses include ethnic identities among the non-UK born population in England and Wales, labour market participation, trends in general health and unpaid care by ethnic group, and inter-ethnic relationships. ",
        "title": "Ethnicity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/ethnicity",
                "id": "ethnicity"
            },
            "content": {
                "href": "http://localhost:25300/topics/ethnicity/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "language",
    "next": {
        "id": "language",
        "description": "The proportion of adults using English as their main language and a breakdown of other languages used in the UK.",
        "title": "Language",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/language",
                "id": "language"
            },
            "content": {
                "href": "http://localhost:25300/topics/language/content"
            }
        }
    },
    "current": {
        "id": "language",
        "description": "The proportion of adults using English as their main language and a breakdown of other languages used in the UK.",
        "title": "Language",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/language",
                "id": "language"
            },
            "content": {
                "href": "http://localhost:25300/topics/language/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "religion",
    "next": {
        "id": "religion",
        "description": "Analyses include people affiliating with a religion in the 2011 census, and religions among the non-UK born population in England and Wales.",
        "title": "Religion",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/religion",
                "id": "religion"
            },
            "content": {
                "href": "http://localhost:25300/topics/religion/content"
            }
        }
    },
    "current": {
        "id": "religion",
        "description": "Analyses include people affiliating with a religion in the 2011 census, and religions among the non-UK born population in England and Wales.",
        "title": "Religion",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/religion",
                "id": "religion"
            },
            "content": {
                "href": "http://localhost:25300/topics/religion/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "sexuality",
    "next": {
        "id": "sexuality",
        "description": "Analyses include sexual identity in the UK by sex, region and age group, sourced from the Annual Population Survey.",
        "title": "Sexual identity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/sexuality",
                "id": "sexuality"
            },
            "content": {
                "href": "http://localhost:25300/topics/sexuality/content"
            }
        }
    },
    "current": {
        "id": "sexuality",
        "description": "Analyses include sexual identity in the UK by sex, region and age group, sourced from the Annual Population Survey.",
        "title": "Sexual identity",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/sexuality",
                "id": "sexuality"
            },
            "content": {
                "href": "http://localhost:25300/topics/sexuality/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "educationandchildcare",
    "next": {
        "id": "educationandchildcare",
        "description": "Early years childcare, school and college education, and higher education and adult learning, including qualifications, personnel, and safety and well-being. ",
        "title": "Education and childcare",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/educationandchildcare",
                "id": "educationandchildcare"
            },
            "content": {
                "href": "http://localhost:25300/topics/educationandchildcare/content"
            }
        }
    },
    "current": {
        "id": "educationandchildcare",
        "description": "Early years childcare, school and college education, and higher education and adult learning, including qualifications, personnel, and safety and well-being. ",
        "title": "Education and childcare",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/educationandchildcare",
                "id": "educationandchildcare"
            },
            "content": {
                "href": "http://localhost:25300/topics/educationandchildcare/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "elections",
    "next": {
        "id": "elections",
        "description": "** no description **",
        "title": "Elections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/elections",
                "id": "elections"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/elections/subtopics"
            }
        },
        "subtopics_ids": [
            "electoralregistration",
            "generalelections",
            "localgovernmentelections"
        ]
    },
    "current": {
        "id": "elections",
        "description": "** no description **",
        "title": "Elections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/elections",
                "id": "elections"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/elections/subtopics"
            }
        },
        "subtopics_ids": [
            "electoralregistration",
            "generalelections",
            "localgovernmentelections"
        ]
    }
})
db.topics.insertOne({
    "id": "electoralregistration",
    "next": {
        "id": "electoralregistration",
        "description": "Annual counts of people listed on electoral registers for the UK and its constituent countries, local government areas and parliamentary constituencies. These figures are based on people registered to vote not people eligible to vote. ",
        "title": "Electoral registration",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/electoralregistration",
                "id": "electoralregistration"
            },
            "content": {
                "href": "http://localhost:25300/topics/electoralregistration/content"
            }
        }
    },
    "current": {
        "id": "electoralregistration",
        "description": "Annual counts of people listed on electoral registers for the UK and its constituent countries, local government areas and parliamentary constituencies. These figures are based on people registered to vote not people eligible to vote. ",
        "title": "Electoral registration",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/electoralregistration",
                "id": "electoralregistration"
            },
            "content": {
                "href": "http://localhost:25300/topics/electoralregistration/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "generalelections",
    "next": {
        "id": "generalelections",
        "description": "Forms of civic engagement, including satisfaction with government and democracy, interest in politics and participation in politics.",
        "title": "General elections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/generalelections",
                "id": "generalelections"
            },
            "content": {
                "href": "http://localhost:25300/topics/generalelections/content"
            }
        }
    },
    "current": {
        "id": "generalelections",
        "description": "Forms of civic engagement, including satisfaction with government and democracy, interest in politics and participation in politics.",
        "title": "General elections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/generalelections",
                "id": "generalelections"
            },
            "content": {
                "href": "http://localhost:25300/topics/generalelections/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "localgovernmentelections",
    "next": {
        "id": "localgovernmentelections",
        "description": "Individuals’ engagement with local government, including involvement in decision-making. ",
        "title": "Local government elections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/localgovernmentelections",
                "id": "localgovernmentelections"
            },
            "content": {
                "href": "http://localhost:25300/topics/localgovernmentelections/content"
            }
        }
    },
    "current": {
        "id": "localgovernmentelections",
        "description": "Individuals’ engagement with local government, including involvement in decision-making. ",
        "title": "Local government elections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/localgovernmentelections",
                "id": "localgovernmentelections"
            },
            "content": {
                "href": "http://localhost:25300/topics/localgovernmentelections/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "healthandsocialcare",
    "next": {
        "id": "healthandsocialcare",
        "description": "Life expectancy and the impact of factors such as occupation, illness and drug misuse. We collect these statistics from registrations and surveys. ",
        "title": "Health and social care",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthandsocialcare",
                "id": "healthandsocialcare"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/healthandsocialcare/subtopics"
            }
        },
        "subtopics_ids": [
            "causesofdeath",
            "childhealth",
            "conditionsanddiseases",
            "disability",
            "drugusealcoholandsmoking",
            "healthcaresystem",
            "healthinequalities",
            "healthandlifeexpectancies",
            "healthandwellbeing",
            "mentalhealth",
            "socialcare"
        ]
    },
    "current": {
        "id": "healthandsocialcare",
        "description": "Life expectancy and the impact of factors such as occupation, illness and drug misuse. We collect these statistics from registrations and surveys. ",
        "title": "Health and social care",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthandsocialcare",
                "id": "healthandsocialcare"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/healthandsocialcare/subtopics"
            }
        },
        "subtopics_ids": [
            "causesofdeath",
            "childhealth",
            "conditionsanddiseases",
            "disability",
            "drugusealcoholandsmoking",
            "healthcaresystem",
            "healthinequalities",
            "healthandlifeexpectancies",
            "healthandwellbeing",
            "mentalhealth",
            "socialcare"
        ]
    }
})
db.topics.insertOne({
    "id": "causesofdeath",
    "next": {
        "id": "causesofdeath",
        "description": "** no description **",
        "title": "Causes of death",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/causesofdeath",
                "id": "causesofdeath"
            },
            "content": {
                "href": "http://localhost:25300/topics/causesofdeath/content"
            }
        }
    },
    "current": {
        "id": "causesofdeath",
        "description": "** no description **",
        "title": "Causes of death",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/causesofdeath",
                "id": "causesofdeath"
            },
            "content": {
                "href": "http://localhost:25300/topics/causesofdeath/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "childhealth",
    "next": {
        "id": "childhealth",
        "description": "** no description **",
        "title": "Child health",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/childhealth",
                "id": "childhealth"
            },
            "content": {
                "href": "http://localhost:25300/topics/childhealth/content"
            }
        }
    },
    "current": {
        "id": "childhealth",
        "description": "** no description **",
        "title": "Child health",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/childhealth",
                "id": "childhealth"
            },
            "content": {
                "href": "http://localhost:25300/topics/childhealth/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "conditionsanddiseases",
    "next": {
        "id": "conditionsanddiseases",
        "description": "Latest data and analysis on coronavirus (COVID-19) in the UK and its effect on the economy and society.",
        "title": "Coronavirus (COVID-19)",
        "keywords": [
            "Coronavirus",
            "COVID 19",
            "corona virus",
            "deaths per day",
            "disease statistics",
            "Covid -19",
            "Covid19",
            "pandemic",
            "epidemic",
            "cover-19",
            "covid19",
            "covid",
            "coroner virus"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/conditionsanddiseases",
                "id": "conditionsanddiseases"
            },
            "content": {
                "href": "http://localhost:25300/topics/conditionsanddiseases/content"
            }
        }
    },
    "current": {
        "id": "conditionsanddiseases",
        "description": "Latest data and analysis on coronavirus (COVID-19) in the UK and its effect on the economy and society.",
        "title": "Coronavirus (COVID-19)",
        "keywords": [
            "Coronavirus",
            "COVID 19",
            "corona virus",
            "deaths per day",
            "disease statistics",
            "Covid -19",
            "Covid19",
            "pandemic",
            "epidemic",
            "cover-19",
            "covid19",
            "covid",
            "coroner virus"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/conditionsanddiseases",
                "id": "conditionsanddiseases"
            },
            "content": {
                "href": "http://localhost:25300/topics/conditionsanddiseases/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "disability",
    "next": {
        "id": "disability",
        "description": "Analysis exploring the lives of disabled people in the UK, to understand disparities and investigate causalities.",
        "title": "Disability",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/disability",
                "id": "disability"
            },
            "content": {
                "href": "http://localhost:25300/topics/disability/content"
            }
        }
    },
    "current": {
        "id": "disability",
        "description": "Analysis exploring the lives of disabled people in the UK, to understand disparities and investigate causalities.",
        "title": "Disability",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/disability",
                "id": "disability"
            },
            "content": {
                "href": "http://localhost:25300/topics/disability/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "drugusealcoholandsmoking",
    "next": {
        "id": "drugusealcoholandsmoking",
        "description": "Smoking and drinking habits in Great Britain, deaths related to drug poisoning  and drug misuse, and deaths caused by diseases known to be related to alcohol consumption.",
        "title": "Drug use, alcohol and smoking",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/drugusealcoholandsmoking",
                "id": "drugusealcoholandsmoking"
            },
            "content": {
                "href": "http://localhost:25300/topics/drugusealcoholandsmoking/content"
            }
        }
    },
    "current": {
        "id": "drugusealcoholandsmoking",
        "description": "Smoking and drinking habits in Great Britain, deaths related to drug poisoning  and drug misuse, and deaths caused by diseases known to be related to alcohol consumption.",
        "title": "Drug use, alcohol and smoking",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/drugusealcoholandsmoking",
                "id": "drugusealcoholandsmoking"
            },
            "content": {
                "href": "http://localhost:25300/topics/drugusealcoholandsmoking/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "healthcaresystem",
    "next": {
        "id": "healthcaresystem",
        "description": "Expenditure on both private and public health care systems in the UK, as well as the results for the National Survey of Bereaved People (VOICES). ",
        "title": "Health care system",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthcaresystem",
                "id": "healthcaresystem"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthcaresystem/content"
            }
        }
    },
    "current": {
        "id": "healthcaresystem",
        "description": "Expenditure on both private and public health care systems in the UK, as well as the results for the National Survey of Bereaved People (VOICES). ",
        "title": "Health care system",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthcaresystem",
                "id": "healthcaresystem"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthcaresystem/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "healthinequalities",
    "next": {
        "id": "healthinequalities",
        "description": "Current patterns and trends in ill health and death by measures of socio-economic status.",
        "title": "Health inequalities",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthinequalities",
                "id": "healthinequalities"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthinequalities/content"
            }
        }
    },
    "current": {
        "id": "healthinequalities",
        "description": "Current patterns and trends in ill health and death by measures of socio-economic status.",
        "title": "Health inequalities",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthinequalities",
                "id": "healthinequalities"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthinequalities/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "healthandlifeexpectancies",
    "next": {
        "id": "healthandlifeexpectancies",
        "description": "Comparisons and inequalities in healthy life expectancies (HLE), disability-free life expectancies (DFLE) and life expectancies (LE).",
        "title": "Health and life expectancies",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthandlifeexpectancies",
                "id": "healthandlifeexpectancies"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthandlifeexpectancies/content"
            }
        }
    },
    "current": {
        "id": "healthandlifeexpectancies",
        "description": "Comparisons and inequalities in healthy life expectancies (HLE), disability-free life expectancies (DFLE) and life expectancies (LE).",
        "title": "Health and life expectancies",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthandlifeexpectancies",
                "id": "healthandlifeexpectancies"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthandlifeexpectancies/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "healthandwellbeing",
    "next": {
        "id": "healthandwellbeing",
        "description": "Analyses of social and economic data from government and other organisations to paint a picture of UK society and how it changes, including comparisons with other countries. ",
        "title": "Health and well-being",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthandwellbeing",
                "id": "healthandwellbeing"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthandwellbeing/content"
            }
        }
    },
    "current": {
        "id": "healthandwellbeing",
        "description": "Analyses of social and economic data from government and other organisations to paint a picture of UK society and how it changes, including comparisons with other countries. ",
        "title": "Health and well-being",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/healthandwellbeing",
                "id": "healthandwellbeing"
            },
            "content": {
                "href": "http://localhost:25300/topics/healthandwellbeing/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "mentalhealth",
    "next": {
        "id": "mentalhealth",
        "description": "People who require social care, as well as those who work in social care, performance levels, and its financial cost. ",
        "title": "Mental health",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/mentalhealth",
                "id": "mentalhealth"
            },
            "content": {
                "href": "http://localhost:25300/topics/mentalhealth/content"
            }
        }
    },
    "current": {
        "id": "mentalhealth",
        "description": "People who require social care, as well as those who work in social care, performance levels, and its financial cost. ",
        "title": "Mental health",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/mentalhealth",
                "id": "mentalhealth"
            },
            "content": {
                "href": "http://localhost:25300/topics/mentalhealth/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "socialcare",
    "next": {
        "id": "socialcare",
        "description": "People who require social care, as well as those who work in social care, performance levels, and its financial cost. ",
        "title": "Social care",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/socialcare",
                "id": "socialcare"
            },
            "content": {
                "href": "http://localhost:25300/topics/socialcare/content"
            }
        }
    },
    "current": {
        "id": "socialcare",
        "description": "People who require social care, as well as those who work in social care, performance levels, and its financial cost. ",
        "title": "Social care",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/socialcare",
                "id": "socialcare"
            },
            "content": {
                "href": "http://localhost:25300/topics/socialcare/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "householdcharacteristics",
    "next": {
        "id": "householdcharacteristics",
        "description": "The composition of households, including those who live alone, overcrowding and under-occupation, as well as internet and social media usage by household.",
        "title": "Household characteristics",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/householdcharacteristics",
                "id": "householdcharacteristics"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/householdcharacteristics/subtopics"
            }
        },
        "subtopics_ids": [
            "homeinternetandsocialmediausage"
        ],
        "spotlight": [
            {
                "href": "/peoplepopulationandcommunity/culturalidentity/sexuality/bulletins/integratedhouseholdsurvey/latest",
                "title": "Integrated Household Survey (Experimental statistics)"
            }
        ]
    },
    "current": {
        "id": "householdcharacteristics",
        "description": "The composition of households, including those who live alone, overcrowding and under-occupation, as well as internet and social media usage by household.",
        "title": "Household characteristics",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/householdcharacteristics",
                "id": "householdcharacteristics"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/householdcharacteristics/subtopics"
            }
        },
        "subtopics_ids": [
            "homeinternetandsocialmediausage"
        ],
        "spotlight": [
            {
                "href": "/peoplepopulationandcommunity/culturalidentity/sexuality/bulletins/integratedhouseholdsurvey/latest",
                "title": "Integrated Household Survey (Experimental statistics)"
            }
        ]
    }
})
db.topics.insertOne({
    "id": "homeinternetandsocialmediausage",
    "next": {
        "id": "homeinternetandsocialmediausage",
        "description": "Use of the internet and social media and how long we use it for. These data include the number of people who have used the internet, those who use the internet frequently and those who have never used the internet. Most of these statistics are broken down by age and gender.",
        "title": "Home internet and social media usage",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/homeinternetandsocialmediausage",
                "id": "homeinternetandsocialmediausage"
            },
            "content": {
                "href": "http://localhost:25300/topics/homeinternetandsocialmediausage/content"
            }
        }
    },
    "current": {
        "id": "homeinternetandsocialmediausage",
        "description": "Use of the internet and social media and how long we use it for. These data include the number of people who have used the internet, those who use the internet frequently and those who have never used the internet. Most of these statistics are broken down by age and gender.",
        "title": "Home internet and social media usage",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/homeinternetandsocialmediausage",
                "id": "homeinternetandsocialmediausage"
            },
            "content": {
                "href": "http://localhost:25300/topics/homeinternetandsocialmediausage/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "housing",
    "next": {
        "id": "housing",
        "description": "Property price, private rent and household survey and census statistics, used by government and other organisations for the creation and fulfilment of housing policy in the UK.",
        "title": "Housing",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/housing",
                "id": "housing"
            },
            "content": {
                "href": "http://localhost:25300/topics/housing/content"
            }
        }
    },
    "current": {
        "id": "housing",
        "description": "Property price, private rent and household survey and census statistics, used by government and other organisations for the creation and fulfilment of housing policy in the UK.",
        "title": "Housing",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/housing",
                "id": "housing"
            },
            "content": {
                "href": "http://localhost:25300/topics/housing/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "leisureandtourism",
    "next": {
        "id": "leisureandtourism",
        "description": "Visits and visitors to the UK, the reasons for visiting and the amount of money they spent here. Also UK residents travelling abroad, their reasons for travel and the amount of money they spent. The statistics on UK residents travelling abroad are an informal indicator of living standards.",
        "title": "Leisure and tourism",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/leisureandtourism",
                "id": "leisureandtourism"
            },
            "content": {
                "href": "http://localhost:25300/topics/leisureandtourism/content"
            }
        }
    },
    "current": {
        "id": "leisureandtourism",
        "description": "Visits and visitors to the UK, the reasons for visiting and the amount of money they spent here. Also UK residents travelling abroad, their reasons for travel and the amount of money they spent. The statistics on UK residents travelling abroad are an informal indicator of living standards.",
        "title": "Leisure and tourism",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/leisureandtourism",
                "id": "leisureandtourism"
            },
            "content": {
                "href": "http://localhost:25300/topics/leisureandtourism/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "personalandhouseholdfinances",
    "next": {
        "id": "personalandhouseholdfinances",
        "description": "Earnings and household spending, including household and personal debt, expenditure, and income and wealth. These statistics help build a picture of our spending and saving decisions. ",
        "title": "Personal and household finances",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/personalandhouseholdfinances",
                "id": "personalandhouseholdfinances"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/personalandhouseholdfinances/subtopics"
            }
        },
        "subtopics_ids": [
            "debt",
            "expenditure",
            "incomeandwealth",
            "pensionssavingsandinvestments"
        ]
    },
    "current": {
        "id": "personalandhouseholdfinances",
        "description": "Earnings and household spending, including household and personal debt, expenditure, and income and wealth. These statistics help build a picture of our spending and saving decisions. ",
        "title": "Personal and household finances",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/personalandhouseholdfinances",
                "id": "personalandhouseholdfinances"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/personalandhouseholdfinances/subtopics"
            }
        },
        "subtopics_ids": [
            "debt",
            "expenditure",
            "incomeandwealth",
            "pensionssavingsandinvestments"
        ]
    }
})
db.topics.insertOne({
    "id": "debt",
    "next": {
        "id": "debt",
        "description": "Debt of UK households, broken down by financial debt and property debt.",
        "title": "Debt",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/debt",
                "id": "debt"
            },
            "content": {
                "href": "http://localhost:25300/topics/debt/content"
            }
        }
    },
    "current": {
        "id": "debt",
        "description": "Debt of UK households, broken down by financial debt and property debt.",
        "title": "Debt",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/debt",
                "id": "debt"
            },
            "content": {
                "href": "http://localhost:25300/topics/debt/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "expenditure",
    "next": {
        "id": "expenditure",
        "description": "Spending patterns of UK households, with findings taken from the Living Costs and Food Survey (LCF).",
        "title": "Expenditure",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/expenditure",
                "id": "expenditure"
            },
            "content": {
                "href": "http://localhost:25300/topics/expenditure/content"
            }
        }
    },
    "current": {
        "id": "expenditure",
        "description": "Spending patterns of UK households, with findings taken from the Living Costs and Food Survey (LCF).",
        "title": "Expenditure",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/expenditure",
                "id": "expenditure"
            },
            "content": {
                "href": "http://localhost:25300/topics/expenditure/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "incomeandwealth",
    "next": {
        "id": "incomeandwealth",
        "description": "UK households income and wealth.",
        "title": "Income and wealth",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/incomeandwealth",
                "id": "incomeandwealth"
            },
            "content": {
                "href": "http://localhost:25300/topics/incomeandwealth/content"
            }
        }
    },
    "current": {
        "id": "incomeandwealth",
        "description": "UK households income and wealth.",
        "title": "Income and wealth",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/incomeandwealth",
                "id": "incomeandwealth"
            },
            "content": {
                "href": "http://localhost:25300/topics/incomeandwealth/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "pensionssavingsandinvestments",
    "next": {
        "id": "pensionssavingsandinvestments",
        "description": "State, private, occupational and workplace pensions, including pension trends that draws together data from other government departments and organisations.  ",
        "title": "Pensions, savings and investments",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/pensionssavingsandinvestments",
                "id": "pensionssavingsandinvestments"
            },
            "content": {
                "href": "http://localhost:25300/topics/pensionssavingsandinvestments/content"
            }
        }
    },
    "current": {
        "id": "pensionssavingsandinvestments",
        "description": "State, private, occupational and workplace pensions, including pension trends that draws together data from other government departments and organisations.  ",
        "title": "Pensions, savings and investments",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/pensionssavingsandinvestments",
                "id": "pensionssavingsandinvestments"
            },
            "content": {
                "href": "http://localhost:25300/topics/pensionssavingsandinvestments/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "populationandmigration",
    "next": {
        "id": "populationandmigration",
        "description": "Size, age, sex and geographic distribution of the UK population, and changes in the UK population and the factors driving these changes. These statistics have a wide range of uses. Central government, local government and the health sector use them for planning, resource allocation and managing the economy. They are also used by people such as market researchers and academics.",
        "title": "Population and migration",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/populationandmigration",
                "id": "populationandmigration"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/populationandmigration/subtopics"
            }
        },
        "subtopics_ids": [
            "internationalmigration",
            "migrationwithintheuk",
            "populationestimates",
            "populationprojections"
        ],
        "spotlight": [
            {
                "href": "/peoplepopulationandcommunity/populationandmigration/populationestimates/articles/overviewoftheukpopulation/latest",
                "title": "Overview of the UK population"
            },
            {
                "href": "/peoplepopulationandcommunity/populationandmigration/internationalmigration/articles/noteonthedifferencebetweennationalinsurancenumberregistrationsandtheestimateoflongterminternationalmigration/latest",
                "title": "Note on the difference between National Insurance number registrations and the estimate of long-term international migration"
            }
        ]
    },
    "current": {
        "id": "populationandmigration",
        "description": "Size, age, sex and geographic distribution of the UK population, and changes in the UK population and the factors driving these changes. These statistics have a wide range of uses. Central government, local government and the health sector use them for planning, resource allocation and managing the economy. They are also used by people such as market researchers and academics.",
        "title": "Population and migration",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/populationandmigration",
                "id": "populationandmigration"
            },
            "subtopics": {
                "href": "http://localhost:25300/topics/populationandmigration/subtopics"
            }
        },
        "subtopics_ids": [
            "internationalmigration",
            "migrationwithintheuk",
            "populationestimates",
            "populationprojections"
        ],
        "spotlight": [
            {
                "href": "/peoplepopulationandcommunity/populationandmigration/populationestimates/articles/overviewoftheukpopulation/latest",
                "title": "Overview of the UK population"
            },
            {
                "href": "/peoplepopulationandcommunity/populationandmigration/internationalmigration/articles/noteonthedifferencebetweennationalinsurancenumberregistrationsandtheestimateoflongterminternationalmigration/latest",
                "title": "Note on the difference between National Insurance number registrations and the estimate of long-term international migration"
            }
        ]
    }
})
db.topics.insertOne({
    "id": "internationalmigration",
    "next": {
        "id": "internationalmigration",
        "description": "People moving into and out of the UK, long term migration, short term migration, and visitor data providing a picture of those entering and leaving the UK, covering all lengths of stay. All data published from the Centre for International Migration ",
        "title": "International migration",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/internationalmigration",
                "id": "internationalmigration"
            },
            "content": {
                "href": "http://localhost:25300/topics/internationalmigration/content"
            }
        }
    },
    "current": {
        "id": "internationalmigration",
        "description": "People moving into and out of the UK, long term migration, short term migration, and visitor data providing a picture of those entering and leaving the UK, covering all lengths of stay. All data published from the Centre for International Migration ",
        "title": "International migration",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/internationalmigration",
                "id": "internationalmigration"
            },
            "content": {
                "href": "http://localhost:25300/topics/internationalmigration/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "migrationwithintheuk",
    "next": {
        "id": "migrationwithintheuk",
        "description": "Residential moves between different geographic areas within the UK sourced from the NHS Patient Register, the NHS Central Register (NHSCR) and the Higher Education Statistics Agency (HESA). ",
        "title": "Migration within the UK",
        "keywords": [
            "internal migration"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/migrationwithintheuk",
                "id": "migrationwithintheuk"
            },
            "content": {
                "href": "http://localhost:25300/topics/migrationwithintheuk/content"
            }
        }
    },
    "current": {
        "id": "migrationwithintheuk",
        "description": "Residential moves between different geographic areas within the UK sourced from the NHS Patient Register, the NHS Central Register (NHSCR) and the Higher Education Statistics Agency (HESA). ",
        "title": "Migration within the UK",
        "keywords": [
            "internal migration"
        ],
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/migrationwithintheuk",
                "id": "migrationwithintheuk"
            },
            "content": {
                "href": "http://localhost:25300/topics/migrationwithintheuk/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "populationestimates",
    "next": {
        "id": "populationestimates",
        "description": "Annual population estimates for the UK and its constituent countries, the regions and counties of England, and local authorities and their equivalents. Estimates for lower and middle layer Super Output Areas, Westminster parliamentary constituencies, electoral wards and National Parks in England and Wales and clinical commissioning groups in England.",
        "title": "Population estimates",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/populationestimates",
                "id": "populationestimates"
            },
            "content": {
                "href": "http://localhost:25300/topics/populationestimates/content"
            }
        }
    },
    "current": {
        "id": "populationestimates",
        "description": "Annual population estimates for the UK and its constituent countries, the regions and counties of England, and local authorities and their equivalents. Estimates for lower and middle layer Super Output Areas, Westminster parliamentary constituencies, electoral wards and National Parks in England and Wales and clinical commissioning groups in England.",
        "title": "Population estimates",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/populationestimates",
                "id": "populationestimates"
            },
            "content": {
                "href": "http://localhost:25300/topics/populationestimates/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "populationprojections",
    "next": {
        "id": "populationprojections",
        "description": "Population projections provide an indication of the future size and age structure of the population based on mid-year population estimates and a set of assumptions of future fertility, mortality and migration. Available for the UK and its constituent countries as national population projections and for the regions, local authorities and clinical commissioning groups in England as subnational population projections. These projections are widely used for resource allocation and planning.",
        "title": "Population projections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/populationprojections",
                "id": "populationprojections"
            },
            "content": {
                "href": "http://localhost:25300/topics/populationprojections/content"
            }
        }
    },
    "current": {
        "id": "populationprojections",
        "description": "Population projections provide an indication of the future size and age structure of the population based on mid-year population estimates and a set of assumptions of future fertility, mortality and migration. Available for the UK and its constituent countries as national population projections and for the regions, local authorities and clinical commissioning groups in England as subnational population projections. These projections are widely used for resource allocation and planning.",
        "title": "Population projections",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/populationprojections",
                "id": "populationprojections"
            },
            "content": {
                "href": "http://localhost:25300/topics/populationprojections/content"
            }
        }
    }
})
db.topics.insertOne({
    "id": "wellbeing",
    "next": {
        "id": "wellbeing",
        "description": "Societal and personal well-being in the UK looking beyond what we produce, to areas such as health, relationships, education and skills, what we do, where we live, our finances and the environment. This data comes from a variety of sources and much of the analysis is new.",
        "title": "Well-being",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/wellbeing",
                "id": "wellbeing"
            },
            "content": {
                "href": "http://localhost:25300/topics/wellbeing/content"
            }
        }
    },
    "current": {
        "id": "wellbeing",
        "description": "Societal and personal well-being in the UK looking beyond what we produce, to areas such as health, relationships, education and skills, what we do, where we live, our finances and the environment. This data comes from a variety of sources and much of the analysis is new.",
        "title": "Well-being",
        "state": "published",
        "links": {
            "self": {
                "href": "http://localhost:25300/topics/wellbeing",
                "id": "wellbeing"
            },
            "content": {
                "href": "http://localhost:25300/topics/wellbeing/content"
            }
        }
    }
})
db.topics.find().forEach(function(doc) {
    printjson(doc);
})

'use strict';

angular.module("tukdesk")
    .constant("apiPrefix", "/apis/v1")
    .constant("errorCode", {
        brandNotFound: 110102
    })
    .constant("ticketStatus", [
        {
            name: "待处理",
            val: "PENDING"
        },
        {
            name: "已回复",
            val: "REPLIED"
        },
        {
            name: "有反馈",
            val: "RESUBMITTED"
        },
        {
            name: "已解决",
            val: "DONE"
        }
    ])
    .constant("ticketPriority", [
        {
            name: "低",
            val: "LOW"
        },
        {
            name: "普通",
            val: "NORMAL"
        },
        {
            name: "优先",
            val: "HIGH"
        },
        {
            name: "紧急",
            val: "URGENT"
        }
    ])
    .constant("ticketFilterOptions", {
        "status": [
            {
                "name": "待处理",
                "val": "PENDING"
            },
            {
                "name": "已回复",
                "val": "REPLIED"
            },
            {
                "name": "有反馈",
                "val": "RESUBMITTED"
            },
            {
                "name": "已解决",
                "val": "DONE"
            }
        ],
        "priority": [
            {
                "name": "低",
                "val": "LOW"
            },
            {
                "name": "普通",
                "val": "NORMAL"
            },
            {
                "name": "优先",
                "val": "HIGH"
            },
            {
                "name": "紧急",
                "val": "URGENT"
            }
        ],
        "time": [
            {
                name: "24小时以内",
                val: "1day"
            },
            {
                name: "48小时以内",
                val: "2days"
            },
            {
                name: "72小时以内",
                val: "3days"
            },
            {
                name: "最近1周",
                val: "7days"
            },
            {
                name: "最近2周",
                val: "14days"
            },
            {
                name: "最近1个月",
                val: "30days"
            }
        ],
        "sort": [
            {
                "name": "更新时间倒序",
                "val": "-updated"
            },
            {
                "name": "更新时间正序",
                "val": "updated"
            },
            {
                "name": "优先级从高到低",
                "val": "-priority"
            },
            {
                "name": "优先级从低到高",
                "val": "priority"
            }
        ]
    })
;

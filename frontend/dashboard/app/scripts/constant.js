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
;

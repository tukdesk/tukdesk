'use strict';

angular.module("tukdesk")
    .factory("utils", function(ticketStatus, ticketPriority, commentTypes) {
        var fac = {};

        fac.ticketStatusClass = function(statusStr) {
            if (ticketStatus.some(function(val) {return val.val === statusStr}) === true) {
                return "tuk-ticket-status-" + statusStr.toLowerCase()
            }
            return ""
        };

        fac.ticketStatusStr = function(statusStr) {
            var str = "";
            ticketStatus.some(function(val) {
                var is = val.val === statusStr;
                if (is === true) {
                    str = val.name;
                }
                return is
            });
            return str
        };

        fac.ticketPriorityClass = function(priorityStr) {
            if (ticketPriority.some(function(val) {return val.val === priorityStr}) === true) {
                return "tuk-ticket-priority-" + priorityStr.toLowerCase()
            }
            return ""
        };

        fac.commentTypeClass = function(typeStr) {
            if (commentTypes.some(function(val) {return val.val === typeStr}) === true) {
                return "tuk-comment-type-" + typeStr.toLowerCase();
            }
            return ""
        };

        fac.commentTypeStr = function(typeStr) {
            var str = "";
            commentTypes.some(function(val) {
                var is = val.val === typeStr;
                if (is === true) {
                    str = val.name
                }
                return is
            });
            return str
        };

        fac.sinceThenStr = function(then) {
            var minute = 60;
            var hour = 60 * minute;
            var day = 24 * hour;
            var month = 30 * day;

            var now = (new Date()).valueOf() / 1000;
            var delta = now - then;
            if (delta > month) {
                var date = new Date(then * 1000);
                var dateStr = date.getMonth() + "-" + date.getDay() + " " + date.getHours() + ":" + date.getMinutes();
                if (date.getFullYear() === now.getFullYear()) {
                    return dateStr
                } else {
                    return date.getFullYear() + "-" + dateStr;
                }
            } else if (delta > day) {
                return parseInt(delta / day, 10) + "天前"
            } else if (delta > hour) {
                return parseInt(delta / hour, 10) + "小时前";
            } else if (delta > minute) {
                return parseInt(delta / minute, 10) + "分钟前";
            } else {
                return "几秒前";
            }
        };

        fac.timestampNow = function() {
            return (new Date()).valueOf() / 1000;
        };

        return fac;
    });

'use strict';

angular.module("tukdesk")
    .factory("utils", function(ticketStatus, ticketPriority) {
        var fac = {};

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

        return fac;
    });

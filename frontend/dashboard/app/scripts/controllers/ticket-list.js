'use strict';

angular.module("tukdesk")
    .controller("ticketListCtrl", ["$scope", "globals", function($scope, globals) {
        var init = function() {
            globals.ticketsGlobalListLoad()
        };

        init();
    }]);

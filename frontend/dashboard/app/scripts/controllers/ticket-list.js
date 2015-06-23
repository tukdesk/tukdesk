'use strict';

angular.module("tukdesk")
    .controller("ticketListCtrl", ["$scope", "globals", "ticketFilterOptions", function($scope, globals, ticketFilterOptions) {
        $scope.ticketFilterOptions = ticketFilterOptions;

        $scope.ticketListRefresh = function() {
            // get comments along with the tickets
            globals.ticketsGlobal.view.filter.comments = "1";
            console.log(globals.ticketsGlobal.view.filter);
            globals.ticketsGlobalListRefresh();
        };

        $scope.showAllId = "";
        $scope.showAllComments = function(ticketId) {
            return ticketId === $scope.showAllId;
        };

        var ticketListLoad = function() {
            globals.ticketsGlobal.view.filter.comments = "1";
            globals.ticketsGlobalListLoad()
        };

        var init = function() {
            ticketListLoad();
        };

        init();
    }]);

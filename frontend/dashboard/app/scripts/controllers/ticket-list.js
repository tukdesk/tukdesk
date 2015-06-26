'use strict';

angular.module("tukdesk")
    .controller("ticketListCtrl", ["$scope", "globals", "api", "ticketFilterOptions", "ticketPriority", "broadcastEvents", function($scope, globals, api, ticketFilterOptions, ticketPriority, broadcastEvents) {
        $scope.ticketFilterOptions = ticketFilterOptions;
        $scope.ticketPriority = ticketPriority;

        $scope.ticketListRefresh = function() {
            // get comments along with the tickets
            globals.ticketsGlobalListRefresh();
        };

        $scope.showAllId = "";
        $scope.showAllComments = function(ticketId) {
            return ticketId === $scope.showAllId;
        };
        $scope.showAllCommentsToggle = function(ticketId) {
            if ($scope.showAllId === ticketId) {
                $scope.showAllId = "";
            } else {
                $scope.showAllId = ticketId;
            }
        };

        var ticketListLoad = function() {
            globals.ticketsGlobalListLoad()
        };

        $scope.ticketUpdate = function(ticket) {
            // 如果需要等到从服务器更新数据后再改变展现, 则要将 copy 动作放到各个具体的 update func
            // 放在此处, 则不论修改是否生效, 展现都发生变化
            var copied = angular.copy(ticket);

            delete copied.comments;
            delete copied.newComment;

            return api.resTickets.update(copied)
                .$promise.then(function(data) {
                    $scope.$emit(broadcastEvents.ticketRefreshWithData, data);
                }, api.responseErr());
        };

        $scope.ticketUpdatePriority = function(ticket, priority) {
            if (ticket.priority === priority) {
                return
            }

            ticket.priority = priority;
            return $scope.ticketUpdate(ticket)
        };

        $scope.ticketSetToSolved = function(ticket) {
            if ($scope.ticketIsSolved(ticket)) {
                return
            }

            ticket.status = "SOLVED";
            return $scope.ticketUpdate(ticket)
        };

        $scope.ticketIsSolved = function(ticket) {
            return ticket.status === "SOLVED"
        };

        $scope.newCommentReset = function(ticket) {
            ticket.newComment = {};
        };

        $scope.commentAddSubmit = function(ticket, commentType) {
            var ticketId = ticket.id;
            ticket.newComment["type"] = commentType;
            api.resComments.add({ticketId: ticket.id}, ticket.newComment)
                .$promise.then(function(data) {
                    ticket.comments.push(data);
                    $scope.newCommentReset(ticket);
                    $scope.$emit(broadcastEvents.ticketRefreshWithId, ticketId);
                }, api.responseErr())
        };

        var init = function() {
            ticketListLoad();
        };

        init();
    }]);

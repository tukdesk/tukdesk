'use strict';

angular.module("tukdesk").run(function($rootScope, $location, api, globals, utils, broadcastEvents) {
    console.log("running");

    api.setDefaultHeader();
    $rootScope.globals = globals;
    $rootScope.utils = utils;
    $rootScope.jumpTo = function(path) {
        $location.path(path);
    };

    $rootScope.$on(broadcastEvents.ticketRefreshWithData, function(e, ticket) {
        console.log("root scope on ticket refresh " + ticket.id);
        globals.ticketsGlobalRefreshWithTicketData(ticket)
    });

    $rootScope.$on(broadcastEvents.ticketRefreshWithId, function(e, ticketId) {
        console.log("root scope on ticket refresh " + ticketId);
        globals.ticketsGlobalRefreshWithTicketId(ticketId)
    });

    $rootScope.$on(broadcastEvents.ticketListRefresh, function(e) {
        console.log("root scope on ticket list refresh");
        globals.ticketsGlobalListLoadedReset();
    });
});

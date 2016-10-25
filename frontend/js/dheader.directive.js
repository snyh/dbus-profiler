(function() {
    'use strick';
    angular.module('dbus-profiler')
        .component('dHeader', {
            templateUrl: "templates/dheader.directive.html",
            controller: ['$scope', 'dapi', function($scope, dapi) {
                var fetchFn = dapi.BuildHeaderInfo()

                var update = function() {
                    fetchFn().then(function(d) {
                        $scope.info = d
                    })
                }
                setInterval(update, 1000);
                update();
            }]
        })
})()

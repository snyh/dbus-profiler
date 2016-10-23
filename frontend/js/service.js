(function() {
    'use strict';
    angular.module('dbus-profiler')
        .factory('dapi', ['$http', function($http) {
            var base = "/dbus/api/"
            return {
                get: get,
            }
            function get(name, cb) {
                var url = base + '/' + name
                $http({method: 'GET', url: url}).then(
                    function(resp) {
                        cb(resp.data)
                    },
                    function(err) {
                        console.log("Errr on get ", url, err)
                    }
                );
            }
        }])
})()

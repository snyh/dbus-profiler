(function() {
    'use strict';
    angular.module('dbus-profiler')
        .factory('dapi', ['$http', function($http) {
            return {
                BuildHeaderInfo: function() {
                    return function() {
                        return Get("/info")
                    }
                },

                BuildIfcInfo: function(name) {
                    return function() {
                        return Get("interface?name="+name)
                    }
                },
                BuildMethodInfo: function(ifc, type, name) {
                    return function() {
                        return MethodInfos(ifc, type, name)
                    }
                },
                IfcInfos: ListIfcInfos
            }
            function Get(name) {
                var base = "/dbus/api/"
                return $http.get(base + '/' + name).then(
                    function(resp) {
                        return resp.data
                    },
                    function(err) {
                        console.log("Errr on get ", url, err)
                    }
                );
            }
            function ListIfcInfos() {
                return Get("/main?1s")
                    .then(function(data) {
                        return data.map(function(d) {
                            return {
                                name: d.Ifc,
                                cost: d.TotalCost
                            }
                        })
                    })
            }
            function MethodInfos(ifcName, type, methodName) {
                return Get("interface?name="+ifcName).then(function(data) {
                    var v
                    switch (type) {
                    case "M":
                        v = data.Method[methodName]
                        break;
                    case "S":
                        v = data.Signal[methodName]
                        break;
                    case "P":
                        v = data.Property[methodName]
                        break;
                    }
                    return {
                        Ifc: ifcName,
                        Method: methodName,
                        Value: {
                            Total: v.Total,
                            Cost: v.Cost.map(function(d) { return d / 1000 / 1000; })
                        }
                    }
                })
            }
        }])
})()

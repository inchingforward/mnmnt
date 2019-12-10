var RecentMemories = RecentMemories || (function() {
    var infoWindow, map, bounds, markers, prevMarker;
    
    mapboxgl.accessToken = 'pk.eyJ1IjoibWphbmdlciIsImEiOiJjazN6NHZlNHkwMjZiM2tudzRpN3FyNzc0In0.avUDs9ardvviib8L8HsMSA';
    
    function init() {
        var center = [-90.419197, 38.677811];

        map = new mapboxgl.Map({
            container: 'map',
            style: 'mapbox://styles/mapbox/light-v10',
            zoom: 13,
            center: center
        });
        
        bounds = new mapboxgl.LngLatBounds();
        markers = {};
    }

    function addMemory(latitude, longitude, memoryId, memoryTitle) {
        var lngLat = [longitude, latitude];

        var popup = new mapboxgl.Popup({closeButton: false})
            .setLngLat(lngLat)
            .setHTML('<a href="/memories/' + memoryId + '">' + memoryTitle + '</a>')
            .setMaxWidth("300px");
        
        var marker = new mapboxgl.Marker()
            .setLngLat(lngLat)
            .addTo(map)
            .setPopup(popup);
        
        bounds.extend(lngLat);
        map.fitBounds(bounds, {padding: 30});

        markers[memoryId] = marker;
    }

    function showInfoWindowForMemory(memoryId) {
        var marker = markers[memoryId];
        marker.togglePopup();

        if (prevMarker) {
            prevMarker.togglePopup();
        }

        prevMarker = marker;
    }

    return {
        showInfoWindowForMemory: showInfoWindowForMemory,
        addMemory: addMemory, 
        init: init
    }
})();

var AddMemory = AddMemory || (function() {
    var map, marker, geocoder;

    var delay = function() {
        var timeout = 0;
        return function(callback, ms) {
            clearTimeout(timeout);
            timeout = setTimeout(callback, ms);
        };
    }();

    function init() {
        map = new mapboxgl.Map({
            container: 'address_search_map',
            style: 'mapbox://styles/mapbox/light-v10',
            center: [-90.1994, 38.6270],
            zoom: 13
        });

        geocoder = new MapboxGeocoder({
            accessToken: mapboxgl.accessToken,
            mapboxgl: mapboxgl,
            placeholder: "1. Search for a place or address" 
        });

        document.getElementById('geocoder').appendChild(geocoder.onAdd(map));
        
        // $("#address_text").keyup(function() {
        //     delay(findAddress, 1500);
        // });

        // $("#address_text").keypress(function (e) {
        //     if ((e.which && e.which == 13) || (e.keyCode && e.keyCode == 13)) {
        //         findAddress();
        //     }
        // });

        // geocoder = new MapboxGeocoder({ accessToken: mapboxgl.accessToken });
        geocoder.on('result', function(result) {
            console.log(result);
            console.log(result.result.center);
        });

        console.log(geocoder);
        console.log(geocoder.query);
    }

    function findAddress() {
        var addressText = $("#address_text").val();
        if (!addressText) {
            alert("Please enter an address to search for.");
            return;
        }

        if (marker) {
            marker.setMap(null);
        }

        console.log(geocoder);
        console.log("gonna query...")
        //geocoder.query(addressText);
        geocoder._geocode(addressText);
        console.log("queried");

        // var latitude = $("#latitude");
        // var longitude = $("#longitude");
        // var submitButton = $("#submit_button");
        // var geocoder = new google.maps.Geocoder();

        // $("#map_spinner").show();


        // geocoder.geocode({address: addressText}, function(results, status) {
        //     $("#map_spinner").hide();

        //     if (status == google.maps.GeocoderStatus.OK) {
        //         map.setCenter(results[0].geometry.location);

        //         marker = new google.maps.Marker({
        //             map: map,
        //             position: results[0].geometry.location
        //         });

        //         latitude.val(results[0].geometry.location.lat());
        //         longitude.val(results[0].geometry.location.lng());

        //         submitButton.attr("disabled", false);
        //     } else {
        //         latitude.val("");
        //         longitude.val("");

        //         submitButton.attr("disabled", "disabled");
        //     }
        // });
    }

    function verifyForm() {
        var latitude = $("#latitude").val();
        var longitude = $("#longitude").val();

        if (!latitude || !longitude) {
            alert("Please search and find an address.");
            $("#address_text").focus();
            return false;
        }

        var titleField = $("#title");
        if (!titleField.val()) {
            alert("Please enter a title.");
            titleField.focus();
            return false;
        }

        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }

        return true;
    }

    return {
        init: init,
        findAddress: findAddress
    }
})();

var EditMemory = AddMemory || (function() {
    var map, marker;

    function init(latitude, longitude) {
        var lngLat = [lng, lat];

        var map = new mapboxgl.Map({
            container: 'edit_memory_map',
            style: 'mapbox://styles/mapbox/light-v10',
            zoom: 13,
            center: lngLat
        });
        
        var marker = new mapboxgl.Marker()
            .setLngLat(lngLat)
            .addTo(map);
    }

    function verifyForm() {
        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }

        return true;
    }

    return {
        init: init,
    }
})();

var MemoryDetails = MemoryDetails || (function() {
    function showMap(lat, lng, title) {
        var lngLat = [lng, lat];

        var map = new mapboxgl.Map({
            container: 'memory_map',
            style: 'mapbox://styles/mapbox/light-v10',
            zoom: 13,
            center: lngLat
        });
        
        var marker = new mapboxgl.Marker()
            .setLngLat(lngLat)
            .addTo(map);
    }

    return {
        showMap: showMap
    }
})();
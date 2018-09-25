
function on_load() {
    get_testrun();
    get_tags();
}


function itemsToTable(json) {
    if (!json) {
        return;
    }
    if (typeof json === "string") {
        json = JSON.parse(json)
    }
    var tbl = d3.select("#table-general");
    tbl.selectAll("*").remove();
    for (i in json) {
        var current_id = json[i].id;
        var row = tbl.append("tr");
        row.append("td").append("a")
            .attr("href", "/api/testrun?id="+current_id).text(json[i].run);
        row.append("td").text(json[i].ts);
    }
}

function itemsToList(json) {
    if (!json) {
        return;
    }
    if (typeof json === "string") {
        json = JSON.parse(json)
    }
    var ul = d3.select("#tags");
    ul.selectAll("*").remove();
    for (i in json) {
        var current = json[i];
        ul.append("li").append("a")
            .attr("href", "#")
            .attr("title", current)
            .text(current)
            .on("click", function(){
                var t = d3.select(this);
                get_testruns_by_tag(t.text());
            });
    }
}

function get_testrun() {
    d3.json("api/testrun")
        .get(function(error, jsonData) {
            itemsToTable(jsonData);
        });
}

function get_tags() {
    d3.json("api/tag")
        .get(function(error, jsonData) {
            itemsToList(jsonData);
        });
}

function get_testruns_by_tag(tag) {
    d3.json("api/testrun?tag="+tag)
        .get(function(error, jsonData) {
            itemsToTable(jsonData);
        });
}

function upload_url() {
    var url = d3.select("#url");
    var status = d3.select("#upload-status");
    //window.alert(url.property("value"));
    status.text("uploading....")
    d3.text("api/upload?url="+url.property("value"))
        .get(function(error, msg) {
            // TODO: check returned status
            var s = d3.select("#upload-status");
            s.text(msg);
            on_load();
        });
}


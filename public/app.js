var req;

function reloadData()
{
    var now = new Date();
    // url = 'liveData?' + now.getTime();
    url = '/directions';

    try {
        req = new XMLHttpRequest();
    } catch (e) {
        try {
            req = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (e) {
            try {
                req = new ActiveXObject("Microsoft.XMLHTTP");
            } catch (oc) {
                alert("No AJAX Support");
                return;
            }
        }
    }

    req.onreadystatechange = processReqChange;
    req.open("GET", url, true);
    req.send(null);
}

function processReqChange()
{
    // If req shows "complete"
    if (req.readyState == 4)
    {
        dataDiv = document.getElementById('currentData');
        dataDiv2 = document.getElementById('container');

        // If "OK"
        if (req.status == 200)
        {
            // Set current data text
            var resText = req.responseText;
            dataDiv.innerHTML = resText;
            resText=resText.split("|");
            // dataDiv2.innerHTML = "<span id=\"question\">What direction?</span>" +
            //     "<div><span>"+resText[0]+"</span><a href=\"\">Vote</a>Forward</div>\n" +
            //     "    <div><span>"+resText[1]+"</span><a href=\"\">Vote</a>Backward</div>\n" +
            //     "    <div><span>"+resText[2]+ "</span><a href=\"\">Vote</a>Turn Left</div>\n" +
            //     "    <div><span>"+resText[3]+"</span><a href=\"\">Vote</a>Turn Right</div>"+
            // "</div>";
            forVal = document.getElementById("forward-val");
            forVal.innerText=resText[0];
            $(forVal).parent().animate({width:(parseInt(resText[0])*30+300)+'px'});
            forVal = document.getElementById("backward-val");
            forVal.innerText=resText[1];
            $(forVal).parent().animate({width:(parseInt(resText[1])*30+300)+'px'});
            forVal = document.getElementById("left-val");
            forVal.innerText=resText[2];
            $(forVal).parent().animate({width:(parseInt(resText[2])*30+300)+'px'});
            forVal = document.getElementById("right-val");
            forVal.innerText=resText[3];
            $(forVal).parent().animate({width:(parseInt(resText[3])*30+300)+'px'});

            // forVal.parent().animate({width:parseInt(resText[0])*100+'px'})
            // $("#container div a").click(function() {
            //     $(this).parent().animate({
            //         width: '+=100px'
            //     }, 500);
            //
            //     $(this).prev().html(parseInt($(this).prev().html()) + 1);
            //     return false;
            // });
            // Start new timer (3 sec)
            timeoutID = setTimeout('reloadData()', 1000);
        }
        else
        {
            // Flag error
            dataDiv.innerHTML = '<p>There was a problem retrieving data: ' + req.statusText + '</p>';
        }
    }
}

$(document).ready(function() {
    $("#container div a").click(function() {
        $(this).parent().animate({
            width: '+=100px'
        }, 500);

        $(this).prev().html(parseInt($(this).prev().html()) + 1);
        return false;
    });
});
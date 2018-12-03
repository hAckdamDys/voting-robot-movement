var req;
var req2;
var reqPost;
const urlVotes = '/directions';
const urlCommand = '/getCommand'
const urlPost = '/postDirection';
const reloadCooldownMS = 1000; //cooldown every X milliseconds
const voteStartWidth = 0;
const voteWidth = 7;

function reloadData()
{
    try {
        req = new XMLHttpRequest();
        req2 = new XMLHttpRequest();
    } catch (e) {
        try {
            req = new ActiveXObject("Msxml2.XMLHTTP");
            req2 = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (e) {
            try {
                req = new ActiveXObject("Microsoft.XMLHTTP");
                req2 = new ActiveXObject("Microsoft.XMLHTTP");
            } catch (oc) {
                alert("No AJAX Support");
                return;
            }
        }
    }

    req.onreadystatechange = updateVotes;
    req.open("GET", urlVotes, true);
    req.send(null);
    req2.onreadystatechange = updateCommand;
    req2.open("GET", urlCommand, true);
    req2.send(null);
}

function updateCommand(){
    if (req2.readyState === XMLHttpRequest.DONE && req2.status===200)
    {
        document.getElementById("last-command").innerText="Last Command: "+req2.responseText;
    }
}

function updateVotes()
{
    // If req shows "complete"
    if (req.readyState === XMLHttpRequest.DONE) {

        if (req.status === 200) {
            // Set current data text
            var resText = req.responseText;
            resText = resText.split("|");
            allCommands=["idle-val","forward-val","backward-val","left-val","right-val"];
            for (i = 0; i < allCommands.length; i++) {
                tmpElem = document.getElementById(allCommands[i]);
                tmpElem.innerText = resText[i];
                $(tmpElem).parent().next().animate({width: (parseInt(resText[i]) * voteWidth + voteStartWidth) + 'px'});
            }
            timeoutID = setTimeout('reloadData()', reloadCooldownMS);
        }
        else {
            // Flag error
            console.log('There was a problem retrieving data: ' + req.statusText);
        }
    }

}


$(document).ready(function() {
    $('.no-zoom').bind('touchend', function(e) {
        e.preventDefault();
        // Add your code here.
        $(this).click();
        // This line still calls the standard click event, in case the user needs to interact with the element that is being clicked on, but still avoids zooming in cases of double clicking.
    })
    $("#container div a").click(function() {
        try {
            reqPost = new XMLHttpRequest();
        } catch (e) {
            try {
                reqPost = new ActiveXObject("Msxml2.XMLHTTP");
            } catch (e) {
                try {
                    reqPost = new ActiveXObject("Microsoft.XMLHTTP");
                } catch (oc) {
                    alert("No AJAX Support");
                    return;
                }
            }
        }
        reqPost.open("POST", urlPost, true);
        reqPost.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
        reqPost.onreadystatechange = function() { // Call a function when the state changes.
            if (this.readyState === XMLHttpRequest.DONE){
                if (this.status === 200) {
                    // reloadData();
                    // don't reload since reloading is often enough
                }
                else{
                    console.log('There was a problem post\'ing data: ' + this.statusText);
                }
            }
        };

        reqPost.send("direction="+$(this).prev()[0].id.split("-")[0]);
        return false;
    });
});

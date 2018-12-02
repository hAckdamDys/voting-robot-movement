var req;
var req2;
var reqPost;
const urlVotes = '/directions';
const urlCommand = '/getCommand'
const urlPost = '/postDirection';
const reloadCooldownMS = 1000; //cooldown every X milliseconds
const voteWidth = 15;

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
        document.getElementById("last-command").innerText=req2.responseText;
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
            tmpElem = document.getElementById("idle-val");
            tmpElem.innerText = resText[0];
            $(tmpElem).parent().animate({width: (parseInt(resText[0]) * voteWidth + 300) + 'px'});
            tmpElem = document.getElementById("forward-val");
            tmpElem.innerText = resText[1];
            $(tmpElem).parent().animate({width: (parseInt(resText[1]) * voteWidth + 300) + 'px'});
            tmpElem = document.getElementById("backward-val");
            tmpElem.innerText = resText[2];
            $(tmpElem).parent().animate({width: (parseInt(resText[2]) * voteWidth + 300) + 'px'});
            tmpElem = document.getElementById("left-val");
            tmpElem.innerText = resText[3];
            $(tmpElem).parent().animate({width: (parseInt(resText[3]) * voteWidth + 300) + 'px'});
            tmpElem = document.getElementById("right-val");
            tmpElem.innerText = resText[4];
            $(tmpElem).parent().animate({width: (parseInt(resText[4]) * voteWidth + 300) + 'px'});
            timeoutID = setTimeout('reloadData()', reloadCooldownMS);
        }
        else {
            // Flag error
            console.log('There was a problem retrieving data: ' + req.statusText);
        }
    }

}


$(document).ready(function() {
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




async function Submit(){
    const status_input = document.getElementById("status");

    var Status = parseInt(status_input.value);

    if (Status>100) status_input.value = "100";
    else if (Status < 0 ) status_input.value = "0";
    const Assignee = parseInt(document.getElementById("description").value)

    try{
        if (!isPost) {
            const response = await fetch(`/tasks`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    Name: document.getElementById("name").value,
                    Status,
                    Description: document.getElementById("description").value,
                    Assignee,
                    ID: TaskID,

                })

            });

            const json = await response.json();

            if (json["success"] === true) {
                document.getElementById("status").innerText = "Saved!"

            } else throw new Error(json["message"]);
        }
        else {
            const response = await fetch(`/tasks`, {
                method:'POST',
                headers: {
                    'Content-Type':'application/json'
                },
                body: JSON.stringify({
                    Name: document.getElementById("name").value,
                    Status,
                    Description: document.getElementById("description").value,
                   Assignee
                })

            });

            const json = await response.json();

            if (json["success"] === true){
                document.getElementById("status").innerText = "Created!"
                window.location.replace("/task/"+json["task"]["ID"]+"/edit")


            } else throw new Error(json["message"]);
        }
    }catch (e) {
        document.getElementById("message").innerText = "Error: " + e.message;
        console.error(e);
    }
}


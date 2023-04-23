async function getJSONData(url) {
  planTable = document.querySelector("table tbody");
  planTable.innerHTML = "";
  fetch(url)
    .then((response) => response.json())
    .then((planListData) => {
      const planList = Array.from(planListData);
      planList.forEach((planData, index) => {
        row = document.createElement("tr");
        element = document.createElement("td");
        element.innerHTML = index + 1;
        row.appendChild(element);
        const plan = Object.values(planData);
        for (let index = 2; index < plan.length; index++) {
          element = document.createElement("td");
          element.innerHTML = plan[index];
          if (index == plan.length - 1) {
            switch (plan[index]) {
              case 0:
                element.innerHTML = "Yapılıyor";
                break;
              case 1:
                element.innerHTML = "Bitti";
                break;
              case 2:
                element.innerHTML = "İptal";
                break;
              default:
                break;
            }
          }
          row.appendChild(element);
        }
        row.innerHTML += `
<td>
    <a href="#editPlanModal" class="edit" data-toggle="modal"><i onclick="editData(${planData.ID})" class="material-icons"
        data-toggle="tooltip" title="Düzenle">&#xE254;</i></a>
    <a href="#deletePlanModal" class="delete" data-toggle="modal"><i onclick="deleteData(${planData.ID})" class="material-icons"
        data-toggle="tooltip" title="Sil">&#xE872;</i></a>
</td> `;
        planTable.appendChild(row);
      });
    });
}
$(document).ready(function () {
  $(".alert").hide();
  getJSONData("/api/plans");

  $('[data-toggle="tooltip"]').tooltip();

  const createPlan = document.getElementById("newPlan");
  createPlan.addEventListener("submit", submitFormCreate);

  function submitFormCreate(event) {
    event.preventDefault();

    const formData = new FormData(createPlan);

    fetch("/api/plans", {
      method: "post",
      body: formData,
    }).then((response) => {
      if (response.ok) {
        closeModals();
        getJSONData("/api/plans");
      } else {
        $("#addAlert").show();
      }
    });
  }

  const deletePlan = document.getElementById("deletePlan");
  deletePlan.addEventListener("submit", submitFormDelete);

  function submitFormDelete(event) {
    event.preventDefault();
    const formData = new FormData(deletePlan);
    for (const [_, value] of formData.entries()) {
      id = value;
    }
    fetch("/api/plans/" + id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({}),
    }).then((response) => {
      if (response.ok) {
        closeModals();
        getJSONData("/api/plans");
      }
    });
  }

  const editPlan = document.getElementById("editPlan");
  editPlan.addEventListener("submit", submitFormEdit);

  function submitFormEdit(event) {
    event.preventDefault();
    const formData = new FormData(editPlan);

    for (const [key, value] of formData.entries()) {
      switch (key) {
        case "ID":
          _id = value;
          break;
        case "Plan":
          _plan = value;
          break;
        case "Date":
          _date = value;
          break;
        case "Start":
          _start = value;
          break;
        case "End":
          _end = value;
          break;
        case "Status":
          _status = value;
          break;
        default:
          break;
      }
    }
    fetch("/api/plans/" + _id, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        Plan: _plan,
        Date: _date,
        Start: _start,
        End: _end,
        Status: parseInt(_status),
      }),
    }).then((response) => {
      if (response.ok) {
        closeModals();
        getJSONData("/api/plans");
      } else {
        $("#editAlert").show();
      }
    });
  }

  const editUser = document.getElementById("editUser");
  editUser.addEventListener("submit", submitFormUser);

  function submitFormUser(event) {
    event.preventDefault();
    const formData = new FormData(editUser);

    for (const [key, value] of formData.entries()) {
      if (key == "Username") {
        _username = value;
      } else {
        _password = value;
      }
    }
    fetch("/api/users", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        Username: _username,
        Password: _password,
      }),
    }).then((response) => {
      if (response.ok) {
        closeModals();
        getJSONData("/api/plans");
      } else {
        $("#userAlert").show();
      }
    });
  }
});

function closeModals() {
  $(".modal").modal("hide");
}

function deleteData(id) {
  $("#deleteData").val(id);
}

function editData(id) {
  $("#editData").val(id);
  fetch("/api/plans/" + id)
    .then((response) => response.json())
    .then((data) => {
      const { Plan, Date, Start, End, Status } = data;
      $("#editPlanTextArea").val(Plan);
      $("#editDate").val(Date);
      $("#editStart").val(Start);
      $("#editEnd").val(End);
      $("#editStatus").val(Status);
    })
    .catch((error) => {
      console.error(error);
    });
}

isWeekly = false;
isMonthly = false;

function currentSort() {
  isWeekly = false;
  isMonthly = false;

  getJSONData("/api/plans");
}

async function weeklySort() {
  isWeekly = true;
  isMonthly = false;

  document.getElementById("pageNumber").innerHTML = 1;
  getJSONData("/api/plans/week/1");
}

function monthlySort() {
  isWeekly = false;
  isMonthly = true;

  document.getElementById("pageNumber").innerHTML = 1;
  getJSONData("/api/plans/month/1");
}

async function changePage(step) {
  const response = await fetch("/api/plans/count");
  const jsonData = await response.json();
  page = document.getElementById("pageNumber");
  if (parseInt(page.innerHTML) + step >= 1) {
    if (isWeekly && parseInt(page.innerHTML) + step <= jsonData["Week"]) {
      page.innerHTML = parseInt(page.innerHTML) + step;
      getJSONData("/api/plans/week/" + parseInt(page.innerHTML));
    } else if (
      isMonthly &&
      parseInt(page.innerHTML) + step <= jsonData["Month"]
    ) {
      page.innerHTML = parseInt(page.innerHTML) + step;
      getJSONData("/api/plans/month/" + parseInt(page.innerHTML));
    }
  }
}
function getCookie(name) {
  return document.cookie
    .split("; ")
    .find((cookie) => cookie.startsWith(name + "="))
    .split("=")[1];
}

function account() {
  uid = getCookie("uid");
  fetch("/api/users/" + uid)
    .then((response) => response.json())
    .then((data) => {
      const { Username, Password } = data;
      $("#studentName").val(Username);
      $("#password").val(Password);
    })
    .catch((error) => {
      console.error(error);
    });
}

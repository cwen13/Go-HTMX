const removeFromDb = (item) => {
  fetch(`/delete?item=${item}`, {method: "Delete"})
    .then((res) => {
      if (res.status == 200){
        window.locaiton.pathname="/";
      }
    });
}

const updateDb = (item) => {
  let input = document.getElementById(item);
  let newItem = input.value;
  fetch(`/update?olditem=${item}&newitem=${newitem}`, {method, "PUT"})
    .then((res) => {
      if (res.status == 200){
        alert("Database updated");
        window.location.pathname="/";
      }
    });
}

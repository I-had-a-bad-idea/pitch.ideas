// Source - https://stackoverflow.com/a/15724300
// Posted by kirlich, modified by community. See post 'Timeline' for change history
// Retrieved 2026-07-03, License - CC BY-SA 4.0
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}




export function loggedIn() {
    return getCookie("logged_in") == "True";
}

export function getUsername() {
    return getCookie("username");
}
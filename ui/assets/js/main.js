const formatter = new Intl.DateTimeFormat("default", {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour12: true,
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
})

const convertMysqlDateToLocal = (str) => {
    const t = str.split(/[- :]/);
    const jsDate = new Date(t[0], t[1] - 1, t[2], t[3], t[4], t[5]);

    return new Date(jsDate.getTime() - (jsDate.getTimezoneOffset() * 60000));
}

window.addEventListener("alpine:init", () => {
    Alpine.data("datetime", () => ({
        prettyDateTime() {
            return formatter.format(convertMysqlDateToLocal(this.$root.getAttribute('datetime')))
        }
    }))
})

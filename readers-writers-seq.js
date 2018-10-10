const ds = []

function reader() {
    const out = ds[Math.floor(ds.length*Math.random())]
    console.log(out)
    return out
}

function writier(num) {
    ds.push(num)
    console.log(num)
}

function main() {
	for (let i=0; i<110; i++) {
        if (i % 10 === 0) {
            writier(i)
        } else {
            reader()
        }
    }
}
main()

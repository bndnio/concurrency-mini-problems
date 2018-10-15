const ds = []

function reader() {
    const out = ds[Math.floor(ds.length*Math.random())]
    return out
}

function writer(num) {
    ds.push(num)
}

function main() {
	for (let i=0; i<1000000; i++) {
        if (i % 10 === 0) {
            writer(i)
        } else {
            reader()
        }
    }
}
main()

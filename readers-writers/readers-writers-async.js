const ds = []

async function reader() {
    return ds[Math.floor(ds.length*Math.random())]
}

async function writer(num) {
    ds.push(num)
}

async function main() {
    const promises = []
	for (let i=0; i<1000000; i++) {
        if (i % 10 === 0) {
            promises.push(writer(i))
        } else {
            promises.push(reader())
        }
    }
    await Promise.all(promises)
}
main()

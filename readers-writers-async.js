const ds = []

async function reader() {
    return ds[Math.floor(ds.length*Math.random())]
}

async function writier(num) {
    ds.push(num)
}

async function main() {
    const promises = []
	for (let i=0; i<110; i++) {
        if (i % 10 === 0) {
            promises.push(writier(i))
        } else {
            promises.push(reader())
        }
    }
    await Promise.all(promises)
}
main()

import express from 'express'

const PORT = process.env.PORT || 3000

const app = express()

app.use(express.json())

app.get('/', async (req, res) => {
    console.log(`[APP] request \`/\` from ${req.socket.remoteAddress}`)
    res.status(200).json({ message: 'It works' })
})

const server = app.listen(PORT, () => {
    console.log(`Started server on :${PORT}`)
})

process.on('SIGTERM', () => {
    console.log(`Received SIGTERM, stopping server`)
    server.close()
})

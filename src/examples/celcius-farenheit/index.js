import express from 'express'
import fetch from 'node-fetch'

// utility function to quit with a message, do not call after initialisation
const die = (message) => {
  console.error(message)
  process.exit(1)
}

const port = process.env.PORT || die("Env var PORT is not set")
const sink = process.env.SINK || die("Env var SINK is not set")

const app = express()

app.use(express.json())

app.get('/status', async (req, res) => {
  return res.json({ 'status': 'ok' })
})

app.post('/input', async (req, res) => {
  const temperature = req.body?.temperature
  if (typeof temperature === "undefined" || temperature === null) {
    console.log(`Bad request: missing temperature on ${JSON.stringify(req.body)}`)
    return res.status(400).json({ 'status': 'Bad request: missing temperature' })
  }

  const farenheit = temperature * (9 / 5) + 32

  fetch(`http://${sink}/input`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      temperature: farenheit
    })
  })

  return res.json({ 'status': 'OK' })
})

app.listen(process.env.PORT, () => {
  console.log(`Starting server on port ${port}`)
})

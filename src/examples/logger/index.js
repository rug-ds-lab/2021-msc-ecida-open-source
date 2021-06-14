import express from 'express'

// utility function to quit with a message, do not call after initialisation
const die = (message) => {
  console.error(message)
  process.exit(1)
}

const appPort = process.env.APP_PORT || die('Environment variable APP_PORT is not set')
const spectatePort = process.env.SPECTATE_PORT || die('Environment variable SPECTATE_PORT is not set')

const app = express()
const viewer = express()

app.use(express.json())
viewer.use(express.json())

let items = []

app.post('/input', (req, res) => {
  const count = items.push(req.body)
  console.log(`[pipeline] received message ${JSON.stringify(req.body)}`)
  return res.json({ 'message': 'OK', 'count': count })
})

viewer.get('/', (req, res) => {
  console.log(`[spectate] showing ${items.length} items`)
  return res.json(JSON.stringify(items))
})

app.listen(appPort, () => {
  console.log(`Started the pipeline server on port ${appPort}`)
})

viewer.listen(spectatePort, () => {
  console.log(`Started the spectate server on port ${spectatePort}`)
})
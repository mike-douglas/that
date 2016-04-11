var config = require('config'),
    express = require('express'),
    proxy = require('express-http-proxy'),
    url = require('url'),
    app = express()

app.use('/collect', proxy(config.get('db.host'), {
    filter: function(req, res) {
        return url.parse(req.url, true).query.key == config.get('key')
    },
    forwardPath: function(req, res) {
        var parsed = url.parse(req.url, true),
            query = parsed.query,
            queryString = []

        delete query['key']

        for (var p in query) {
            if (query.hasOwnProperty(p)) {
                queryString.push(`${encodeURIComponent(p)}=${encodeURIComponent(query[p])}`)
            }
        }

        return `${parsed.pathname}?${queryString.join('&')}`
    }
}))

app.use(function(req, res, next) {
    console.log(`${new Date()}: ${req.originalUrl}`)
    next()
})

var server = app.listen(config.get('web.port'), function() {
    var port = server.address().port

    console.log(`Listening on ${port}`)
})

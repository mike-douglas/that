var config = require('config'),
    express = require('express'),
    proxy = require('express-http-proxy'),
    url = require('url'),
    app = express()

app.use('/collect', proxy(config.get('db.host'), {
    filter: function(req, res) {
        return url.parse(req.url, true).query.key == config.get('key') || req.method == 'GET'
    },
    forwardPath: function(req, res) {
        var parsed = url.parse(req.url, true),
            query = parsed.query,
            queryString = []

        delete query['key']
        delete query['callback']

        for (var p in query) {
            if (query.hasOwnProperty(p)) {
                queryString.push(`${encodeURIComponent(p)}=${encodeURIComponent(query[p])}`)
            }
        }

        return `${parsed.pathname}?${queryString.join('&')}`
    },
    intercept: function(rsp, data, req, res, callback) {
        var parsed = url.parse(req.url, true)

        res.setHeader('Content-Type', 'application/json')

        if (parsed.query.hasOwnProperty('callback')) {
            callback(null, `${parsed.query.callback}(${data.toString('utf8')})`)
        } else {
            callback(null, data.toString('utf8')) // JSON.stringify(data))
        }
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

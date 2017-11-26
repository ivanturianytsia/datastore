class Utils {
  constructor (ctx) {
    this.ctx = ctx
  }
  handleErr (err) {
    const h = this.ctx.$createElement
    this.ctx.$message({
      message: h('p', { color: 'red' }, `Oops, error occured: ${err.bodyText}`)
    })
    console.log(err)
    // this.ctx.logout()
  }
}

export default Utils

"use strict";

class UploadFile extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (<div>
            <div className="custom-file">
                <input type="file" className="custom-file-input" id="customFile" />
                <label className="custom-file-label" htmlFor="customFile">选择镜像文件</label>
            </div>
            <div className="custom-file">
                <div>item: xxx</div>
                <div>item: xxx</div>
                <div>item: xxx</div>
                <div>item: xxx</div>
            </div>
            <div className="position-absolute nlv" onClick={this.props.nextPage}>
                <div className={"float-right"}>
                                <span>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                                         fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                         stroke-linejoin="round" className="feather feather-chevrons-right">
                                        <polyline points="13 17 18 12 13 7"></polyline>
                                        <polyline points="6 17 11 12 6 7"></polyline>
                                    </svg>
                                </span>
                    <span>下一步</span>
                </div>
            </div>
        </div>)
    }
}

class PerfectInfo extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (<div>
            <div className="custom-file">
                <form className="col-lg-8">
                    <div className="form-group row">
                        <label htmlFor="staticEmail" className="col-sm-2 col-form-label">FID</label>
                        <div className="col-sm-10">
                            <input type="text" readOnly className="form-control-plaintext" id="staticEmail"
                                   value="email@example.com" />
                        </div>
                    </div>
                    <div className="form-group row">
                        <label htmlFor="inputPassword" className="col-sm-2 col-form-label">Name</label>
                        <div className="col-sm-10">
                            <input type="text" className="form-control" id="inputPassword"
                                   placeholder="Password" />
                        </div>
                    </div>
                    <div className="form-group row">
                        <label htmlFor="inputPassword" className="col-sm-2 col-form-label">Group</label>
                        <div className="col-sm-10">
                            <select className="form-control" id="exampleFormControlSelect1">
                                <option>1</option>
                                <option>2</option>
                                <option>3</option>
                                <option>4</option>
                                <option>5</option>
                            </select>
                        </div>
                    </div>
                    <div className="form-group row">
                        <label htmlFor="inputPassword" className="col-sm-2 col-form-label">Version</label>
                        <div className="col-sm-10">
                            <input type="text" className="form-control" id="inputPassword"
                                   placeholder="Password" />
                        </div>
                    </div>
                    <div className="form-group row">
                        <label htmlFor="inputPassword" className="col-sm-2 col-form-label">Remarks</label>
                        <div className="col-sm-10">
                            <textarea className="form-control" id="exampleFormControlTextarea1"
                                      rows="3"></textarea>
                        </div>
                    </div>
                </form>
            </div>
            <div className="position-absolute nlv">
                <div className={"float-right"} onClick={this.props.nextPage}>
                    <span>
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                             stroke-linejoin="round" className="feather feather-chevrons-right">
                            <polyline points="13 17 18 12 13 7"></polyline>
                            <polyline points="6 17 11 12 6 7"></polyline>
                        </svg>
                    </span>
                    <span>下一步</span>
                </div>
                <div className={"float-right"} onClick={this.props.lastPage}>
                    <span>
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                             stroke-linejoin="round" className="feather feather-chevrons-left">
                            <polyline points="11 17 6 12 11 7"></polyline>
                            <polyline points="18 17 13 12 18 7"></polyline>
                        </svg>
                    </span>
                    <span>上一步</span>
                </div>
            </div>
        </div>)
    }
}

class Msg extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (<div>
            <div className="msg title">
                Oops!!
            </div>
            <div className="msg">
                镜像创建成功! <button className="btn btn-light" onClick={() => window.history.back(-1)}>关闭</button>
            </div>
        </div>)
    }
}

class Main extends React.Component {
    constructor(props) {
        super(props);
        this.state = {active: 0};
        this.nextPage = this.nextPage.bind(this);
        this.lastPage = this.lastPage.bind(this);
    }

    changePage(active = 0) {
        let ops = {active: active};
        switch (active) {
            case 0:
                Object.assign(ops, {name: "上传镜像文件", value: "35"});
                break;
            case 1:
                Object.assign(ops, {name: "填写镜像信息", value: "70"});
                break;
            case 2:
                Object.assign(ops, {name: "通知", value: "100"});
                break;
            default:
                Object.assign(ops, {name: "上传镜像文件", value: "35"});
                break;
        }
        return ops;
    }

    nextPage() {
        const {active = 0} = this.state;
        this.setState(this.changePage(active + 1))
    }

    lastPage() {
        const {active = 0} = this.state;
        this.setState(this.changePage(active - 1))
    }

    switchPage() {
        const {active = 0} = this.state;
        switch (active) {
            case 0:
                return <UploadFile nextPage={this.nextPage}/>;
            case 1:
                return <PerfectInfo nextPage={this.nextPage} lastPage={this.lastPage}/>;
            case 2:
                return <Msg/>;
            default:
                return <UploadFile nextPage={this.nextPage}/>;
        }
    }

    returnLastPage() {
        window.history.back(-1);
    }

    render() {
        const {name = "上传镜像文件", value = "35"} = this.state;
        return (
            <div class="row">
                <main role="main" className="col-lg-8 m-auto px-4">
                    <div
                        className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                        <h1 className="h2">
                            <span onClick={this.returnLastPage}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                                     fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                     stroke-linejoin="round" className="feather feather-corner-down-left">
                                    <polyline points="9 10 4 15 9 20"></polyline><path d="M20 4v7a4 4 0 0 1-4 4H4"></path>
                                </svg>
                            </span>
                            <span>{name}</span>
                        </h1>
                    </div>
                    <div className="my-3 p-3 bg-white rounded shadow-sm m-cent">
                        <div className="progress">
                            <div className="progress-bar progress-bar-striped progress-bar-animated" role="progressbar"
                                 aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style={{"width": value + "%"}}></div>
                        </div>
                        {this.switchPage()}
                    </div>
                </main>
            </div>
        )
    }
}

Main.defaultProps = {
    menu: [],
    links: []
};

ReactDOM.render(<Main />, document.getElementById("app"));
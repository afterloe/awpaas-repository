"use strict";

class UploadFile extends React.Component {
    constructor(props) {
        super(props);
        this.selectFile = this.selectFile.bind(this);
        this.state = {}
    }

    selectFile(event) {
        const content = event.currentTarget;
        if (0 != content.files.length) {
            let file = content.files[0];
            this.setState({uploadFile: file, topic: `已选择：${file.name}`});
        }
    }

    loadSelectFileInfo({lastModified, name = "", size, type}) {
        if ("" === name) {
            return (
                <div className="custom-file show">
                    <p className={"text-primary"}>请选择镜像文件</p>
                </div>
            )
        }
        return (
            <div className="custom-file show">
                <div className="info">
                    <div className={"row"}>
                        <span className={"col"}>文件名:</span>
                        <span className={"col val"}>{name}</span>
                    </div>
                    <div className={"row"}>
                        <span className={"col"}>文件大小:</span>
                        <span className={"col val"}>{size / 1000} kb</span>
                    </div>
                    <div className={"row"}>
                        <span className={"col"}>文件类型:</span>
                        <span className={"col val"}>{type}</span>
                    </div>
                    <div className={"row"}>
                        <span className={"col"}>修改日期:</span>
                        <span className={"col val"}>{getLocalTime(lastModified)}</span>
                    </div>
                </div>
            </div>
        )
    }

    render() {
        const {uploadFile = {}, topic = "选择镜像文件"} = this.state;
        return (
            <div className={"view col-7"}>
                <div className="custom-file">
                    <input type="file" className="custom-file-input" onChange={this.selectFile} />
                    <label className="custom-file-label" htmlFor="customFile">{topic}</label>
                </div>
                {this.loadSelectFileInfo(uploadFile)}
                <div className={"custom-file"}>
                    <button type="button" className="btn btn-primary btn-lg btn-block" onClick={this.props.nextPage}>
                        下一步
                    </button>
                </div>
            </div>)
    }
}

class PerfectInfo extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div className={"view col-7"}>
                <form className="custom-file show">
                    <div className="form-group row">
                        <label htmlFor="staticEmail" className="col-sm-2 col-form-label">FID</label>
                        <div className="col-sm-10">
                            <input type="text" readOnly className="form-control-plaintext" value="email@example.com" />
                        </div>
                    </div>
                    <div className="form-group row">
                        <label htmlFor="inputPassword" className="col-sm-2 col-form-label">Name</label>
                        <div className="col-sm-10">
                            <input type="text" className="form-control" placeholder="image name" />
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
                            <input type="text" className="form-control" placeholder="image version" />
                        </div>
                    </div>
                    <div className="form-group row">
                        <label htmlFor="inputPassword" className="col-sm-2 col-form-label">Remarks</label>
                        <div className="col-sm-10">
                            <textarea className="form-control" rows="3"></textarea>
                        </div>
                    </div>
                </form>
                <div className={"custom-file"}>
                    <button type="button" className="btn btn-light btn-lg btn-block" onClick={this.props.lastPage}>
                        上一步
                    </button>
                    <button type="button" className="btn btn-primary btn-lg btn-block" onClick={this.props.nextPage}>
                        下一步
                    </button>

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
                镜像创建成功! <br />
                <button className="btn btn-light" onClick={() => window.history.back(-1)}>关闭</button>
            </div>
        </div>)
    }
}

class Main extends React.Component {
    constructor(props) {
        super(props);
        this.state = {active: 1};
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
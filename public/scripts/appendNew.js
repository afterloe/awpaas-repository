"use strict";

class Main extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div class="row">
                <main role="main" className="col-lg-8 m-auto px-4">
                    <div
                        className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                        <h1 className="h2">
                            <span>
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                                       fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                       stroke-linejoin="round" className="feather feather-corner-down-left">
                                    <polyline points="9 10 4 15 9 20"></polyline><path d="M20 4v7a4 4 0 0 1-4 4H4"></path>
                                </svg>
                            </span>
                            <span>创建镜像</span>
                        </h1>
                    </div>
                    <div className="my-3 p-3 bg-white rounded shadow-sm m-cent">
                        <div className="progress">
                            <div className="progress-bar progress-bar-striped progress-bar-animated" role="progressbar"
                                 aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style={{"width": "35%"}}></div>
                        </div>
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
                        <div className="position-absolute nlv">
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
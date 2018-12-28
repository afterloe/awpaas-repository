"use strict";

const labels = ["周天", "周一", "周二", "周三", "周四", "周五", "周六"];
const requestTotal = [23, 233, 122, 309, 177, 133, 12];
const cpuTotal = [3, 12, 20, 12, 45, 124, 18];
const reqRank = [{name: "docker", url: "user/images/json", count: "50002", trend: "-"},
    {name: "docker", url: "user/images/json", count: "50002", trend: "-"},
    {name: "docker", url: "user/images/json", count: "50002", trend: "UP"},
    {name: "docker", url: "user/images/json", count: "50002", trend: "DOWN"},
    {name: "docker", url: "user/images/json", count: "50002", trend: "DOWN"},
    {name: "docker", url: "user/images/json", count: "50002", trend: "-"}];

class TotalMain extends React.Component {
    constructor(props) {
        super(props);
        this.state = {}; // 初始化数据
    }

    componentDidMount() {
        // 拉取数据
        TotalMain.loadChart(requestTotal, cpuTotal); // 绘制流量报表
        this.setState({rank: reqRank})
    }

    renderRequestRank() {
        const {rank = []} = this.state;
        return rank.map((it, i) => (
            <tr>
                <td>{i}</td>
                <td>{it.name}</td>
                <td>{it.url}</td>
                <td>{it.count}</td>
                <td>{it.trend}</td>
            </tr>
        ));
    }

    static loadChart(requestTotal, cpuTotal) {
        new Chart(document.getElementById("myChart-1"), {
            type: 'doughnut',
            data: {
                datasets: [{
                    data: [10, 12, 17, 21, 13],
                    backgroundColor: ["#f94d00", "#ffc107", "#007bff", "#28a745", "#6f42c1"]
                }],
                labels: [
                    'awpaas',
                    'public',
                    'spring_boot',
                    'go',
                    '.net'
                ]
            }
        });
        new Chart(document.getElementById("myChart-2"), {
            type: 'horizontalBar',
            data: {
                datasets: [{
                    data: [11, 12, 17, 21, 13],
                    backgroundColor: ["#f94d00", "#ffc107", "#007bff", "#28a745", "#6f42c1"]
                }],
                labels: [
                    'awpaas',
                    'public',
                    'spring_boot',
                    'go',
                    '.net'
                ]
            },
            options: {
                scales: {
                    xAxes: [{
                        stacked:true
                    }],
                    yAxes: [{
                        stacked:true
                    }]
                },
                legend: {
                    display: false,
                }
            }
        });
        new Chart(document.getElementById("myChart-3"), {
            type: 'polarArea',
            data: {
                datasets: [{
                    data: [11, 12, 17, 21, 13],
                    backgroundColor: ["#f94d00", "#ffc107", "#007bff", "#28a745", "#6f42c1"]
                }],
                labels: [
                    'awpaas',
                    'public',
                    'spring_boot',
                    'go',
                    '.net'
                ]
            },
        });
    }

    render() {
        return (
            <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-4">
                <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                    <h1 class="h2">镜像统计</h1>
                </div>
                <div className={"row"}>
                    <canvas className="col-md-6" id="myChart-1" width="1053" height="444"></canvas>
                    <div className={"col-md-6"}>
                        <div className="table-responsive">
                            <table className="table table-striped">
                                <thead>
                                <tr>
                                    <th>#</th>
                                    <th>服务名</th>
                                    <th>URL</th>
                                    <th>访问次数</th>
                                    <th>趋势</th>
                                </tr>
                                </thead>
                                <tbody>{this.renderRequestRank()}</tbody>
                            </table>
                        </div>
                    </div>
                </div>
                <div
                    className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                    <h1 className="h2">更新趋势</h1>
                </div>
                <div className={"row"}>
                    <div className={"col-md-6"}>
                        <div className="table-responsive">
                            <table className="table table-borderless table-hover">
                                <thead className={"thead-light"}>
                                <tr>
                                    <th>#</th>
                                    <th>服务名</th>
                                    <th>URL</th>
                                    <th>访问次数</th>
                                    <th>趋势</th>
                                </tr>
                                </thead>
                                <tbody>{this.renderRequestRank()}</tbody>
                            </table>
                        </div>
                    </div>
                    <canvas className="col-md-6" id="myChart-2" width="1053" height="444"></canvas>
                </div>
                <div
                    className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                    <h1 className="h2">下载排行</h1>
                </div>
                <div className={"row"}>
                    <canvas className="col-md-6" id="myChart-3" width="1053" height="444"></canvas>
                    <div className={"col-md-6"}>
                        <div className="table-responsive">
                            <table className="table table-striped">
                                <thead>
                                <tr>
                                    <th>#</th>
                                    <th>服务名</th>
                                    <th>URL</th>
                                    <th>访问次数</th>
                                    <th>趋势</th>
                                </tr>
                                </thead>
                                <tbody>{this.renderRequestRank()}</tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </main>
        )
    }
}
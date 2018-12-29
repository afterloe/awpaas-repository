"use strict";

const systemMenu = [{
    name: "首页",
    icon: "images/home.svg",
    isClick: true,
    index: "main"
}, {
    name: "镜像仓库",
    icon: "images/layers.svg",
    index: "warehouse"
}, {
    name: "远程仓库",
    icon: "images/cloud.svg",
    index: "serviceRegistry"
}, {
    name: "活跃用户",
    icon: "images/users.svg",
    index: "activeUsers"
}, {
    name: "仓库状态",
    icon: "images/bar-chart-2.svg",
    index: "gatewayStatus"
}];
const linkMenu = [{name: "统一管理子系统", href: "https://127.0.0.1:8088"}, {name: "统一认证及审计子系统", href: "https://127.0.0.1:8088"}, {name: "蜂窝式数据仓库管理子系统", href: "https://127.0.0.1:8088"}, {name: "数据采集清洗引擎", href: "https://127.0.0.1:8088" }];

class Main extends React.Component {
    constructor(props) {
        super(props);
        const {menu = []}= props;
        // TODO
        this.state = {menu, active: "warehouse"};
        this.clickItem = this.clickItem.bind(this);
    }

    clickItem(event) {
        const key = event.currentTarget.getAttribute("data-index") || "";
        if ("" === key) return;
        this.setState(prevState => {
            const {menu} = prevState;
            return {
                active: key,
                menu: menu.map(it => {
                it.isClick = key === it.index;
                return it;
            })}
        });
    }

    switchPage() {
        const {active = "main"} = this.state;
        switch (active) {
            case "main":
                return <TotalMain />;
            case "warehouse":
                return <Warehouse />;
            case "serviceRegistry":
                return <ServiceRegistry />;
            case "activeUsers":
                return <ActiveUsers />;
            case "gatewayStatus":
                return "gatewayStatus";
        }
    }

    render() {
        const {menu} = this.state;
        return (
            <div class="row">
                <NavLeft menu={menu} links={this.props.links} clickItem={this.clickItem}/>
                {this.switchPage()}
            </div>
        )
    }
}

Main.defaultProps = {
    menu: [],
    links: []
};

ReactDOM.render(<Header name="镜像管理" version="v1.0.0"/>, document.getElementById("head"));
ReactDOM.render(<Main menu={systemMenu} links={linkMenu}/>, document.getElementById("app"));
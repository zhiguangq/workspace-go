<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{config "String" "globaltitle" ""}}</title>
{{template "inc/meta.tpl" .}}
</head><body class="sticky-header">
<section> 
  <!-- main content start-->
  <div class="main-content" >
    <!-- header section start-->
    <div class="header-section">
      <!--toggle button start-->
      <a class="toggle-btn"><i class="fa fa-bars"></i></a> 
	
     <div class="menu-right">
       <ul class="notification-menu">
        <li> <a href="javascript:;" class="btn btn-default dropdown-toggle" data-toggle="dropdown"> <img src="{{getAvatar .LoginAvatar}}" alt="{{.LoginUsername}}" /> {{.LoginUsername}} <span class="caret"></span> </a>
        <ul class="dropdown-menu dropdown-menu-usermenu pull-right">
         <li><a href="/logout"><i class="fa fa-sign-out"></i> 退出</a></li>
        </ul>
        </li>
       </ul>
      </div>
	
	</div>
    <!-- header section end-->
    <!-- page heading start-->
    <div class="page-heading">
      <h3> 摄像头管理 </h3>
      <ul class="breadcrumb pull-left">
        <li> <a href="/camera/manage">管理主页</a> </li>
      </ul>
      <div class="pull-right"><a href="/camera/add" class="btn btn-success">添加摄像头</a></div>
    </div>
    <!-- page heading end-->
    <!--body wrapper start-->
    <div class="wrapper">
      <div class="row">
        <div class="col-lg-12">
          <section class="panel">
            <header class="panel-heading"> {{.title}} </header>
            <div class="panel-body">
              <form class="form-horizontal adminex-form" id="camera-form">
                <div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">摄像头名称</label>
                  <div class="col-sm-10">
                    <input type="text" name="name" value="{{.dep.Name}}" class="form-control" placeholder="请填写摄像头名称">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">录像机域名</label>
                  <div class="col-sm-10">
                    <input type="text" name="dns" value="{{.dep.Dns}}" class="form-control" placeholder="请填写域名">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">录像机用户名</label>
                  <div class="col-sm-10">
                    <input type="text" name="users" value="{{.dep.Users}}" class="form-control" placeholder="请填写录像机用户名">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">录像机密码</label>
                  <div class="col-sm-10">
                    <input type="text" name="pass" value="{{.dep.Pass}}" class="form-control" placeholder="请填写录像机密码">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">摄像头通道</label>
                  <div class="col-sm-10">
                    <input type="text" name="channel" value="{{.dep.Channel}}" class="form-control" placeholder="请填写摄像头通道">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">摄像头位置1</label>
                  <div class="col-sm-10">
                    <input type="text" name="address1" value="{{.dep.Address1}}" class="form-control" placeholder="请填写摄像头位置1">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">摄像头位置2</label>
                  <div class="col-sm-10">
                    <input type="text" name="address2" value="{{.dep.Address2}}" class="form-control" placeholder="请填写摄像头位置2">
                  </div>
                </div>
				<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">播放地址</label>
                  <div class="col-sm-10">
                    <input type="text" name="playurl" value="{{.dep.Playurl}}" readonly="readonly" class="form-control" placeholder="不用填写，会自动生成">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label">描述</label>
                  <div class="col-sm-10">
                    <textarea name="desc" placeholder="请填写描述" style="height:90px;" class="form-control">{{.dep.Desc}}</textarea>
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-lg-2 col-sm-2 control-label"></label>
                  <div class="col-lg-10">
                    <input type="hidden" name="id" value="{{.dep.Id}}">
                    <button type="submit" class="btn btn-primary">提 交</button>
                  </div>
                </div>
              </form>
            </div>
          </section>
        </div>
      </div>
    </div>
    <!--body wrapper end-->
    <!--footer section start-->
    <!--footer section end-->
  </div>
  <!-- main content end-->
</section>
{{template "inc/foot.tpl" .}}
</body>
</html>

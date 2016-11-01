<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{config "String" "globaltitle" ""}}</title>
{{template "inc/meta.tpl" .}}
<link href="/static/css/table-responsive.css" rel="stylesheet">
</head><body class="sticky-header">
<section> 
  <!-- main content start-->
  <div class="main-content" >
    <!-- header section start-->
    <div class="header-section">
      <!--toggle button start-->
      <a class="toggle-btn"><i class="fa fa-bars"></i></a>
      <!--toggle button end-->
      <!--search start-->
      <form class="searchform" action="/camera/manage" method="get">
        <select name="status" class="form-control">
          <option value="">状态</option>
          <option value="1" {{if eq "1" .condArr.status}}selected{{end}}>正常</option>
          <option value="2" {{if eq "2" .condArr.status}}selected{{end}}>屏蔽</option>
        </select>
        <input type="text" class="form-control" name="keywords" placeholder="请输入名称" value="{{.condArr.keywords}}"/>
        <button type="submit" class="btn btn-primary">搜索</button>
      </form>
      <!--search end-->
	
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
        <li class="active"> 摄像头 </li>
      </ul>
      <div class="pull-right"><a href="/camera/add" class="btn btn-success">添加摄像头</a></div>
    </div>
    <!-- page heading end-->
    <!--body wrapper start-->
    <div class="wrapper">
      <div class="row">
        <div class="col-sm-12">
          <section class="panel">
            <header class="panel-heading"> 总数：{{.countDepart}}<span class="tools pull-right"><a href="javascript:;" class="fa fa-chevron-down"></a>
              <!--a href="javascript:;" class="fa fa-times"></a-->
              </span> </header>
            <div class="panel-body">
              <section id="unseen">
                <form id="department-form-list">
                  <table class="table table-bordered table-striped table-condensed">
                    <thead>
                      <tr>
                        <th>名称</th>
                        <th>描述</th>
                        <th>状态</th>
                        <th>操作</th>
                      </tr>
                    </thead>
                    <tbody>
                    
                    {{range $k,$v := .depart}}
                    <tr>
                      <td>{{$v.Name}}</td>
                      <td>{{$v.Desc}}</td>
                      <td>{{if eq 1 $v.Status}}正常{{else}}屏蔽{{end}}</td>
                      <td><div class="btn-group">
                          <button type="button" class="btn btn-primary dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false"> 操作<span class="caret"></span> </button>
                          <ul class="dropdown-menu">
                            <li><a href="/camera/edit/{{$v.Id}}">编辑</a></li>
                            <li role="separator" class="divider"></li>
                            {{if eq 1 $v.Status}}
                            <li><a href="javascript:;" class="js-camera-single" data-id="{{$v.Id}}" data-status="2">屏蔽</a></li>
                            {{else}}
                            <li><a href="javascript:;" class="js-camera-single" data-id="{{$v.Id}}" data-status="1">正常</a></li>
                            {{end}}
                          </ul>
                        </div></td>
                    </tr>
                    {{end}}
                    </tbody>
                    
                  </table>
                </form>
                {{template "inc/page.tpl" .}}
				 </section>
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

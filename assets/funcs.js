$(document).ready(function() {
  getInfo()
$('body').on('click', '#getInfoBtn', function(event) {
  event.preventDefault();
  /* Act on the event */
  // TODO ajax request and change all 4 button values
  getInfo()
});

  // TODO power on click
$('body').on('click', '#powerOnBtn', function(event) {
  event.preventDefault();
  /* Act on the event */
  // TODO ajax request and change all 4 button values
  poweron();
  $('.grilltemp').toggleClass('yft');
});

// TODO power off click
$('body').on('click', '#powerOffBtn', function(event) {
  event.preventDefault();
  /* Act on the event */
  // TODO ajax request and change all 4 button values
  poweroff();
  $('.grilltemp').toggleClass('yft');
});

// TODO start logging click
$('#logForm').submit(function(event) {
  event.preventDefault();
  startLog();
});
$('#setGrillForm').submit(function(event) {
  event.preventDefault();
  setgrilltarget();
});
$('#setProbeForm').submit(function(event) {
  event.preventDefault();
  setprobetarget();
});
// TODO start logging click
$('body').on('click', '#broadcastClientBtn', function(event) {
  event.preventDefault();
  /* Act on the event */
  // TODO ajax request and change all 4 button values
  btoc()
  $('.grilltemp').toggleClass('yft');
});
// TODO show history
$('body').on('click', '#historyBtn', function(event) {
  event.preventDefault();
  /* Act on the event */
  // TODO ajax request and change all 4 button values
  history()
  $('.grilltemp').toggleClass('yft');
});

}); // end of onready

/////////////
//
//  FUNCTIONS
//
/////////////

function startLog() {
  $.ajax({
    url: 'log',
    type: 'POST',
    dataType: 'json',
    data: '{"food": "'+$("#foodname").val()+
    '", "weight": '+$("#weight").val()+
    ', "interval": '+$("#interval").val()+'}'
  })
  .done(function(data) {
    console.log("done");
  })
  .fail(function() {
    console.log("error");
  })
  .always(function() {
    console.log("complete");
  });

}
function getInfo() {
$.ajax({
  url: 'info',
  type: 'GET',
  dataType: 'json',
  // data: {param1: 'value1'}
})
.done(function(data) {
  console.log("success");
  console.log(data);
  $('.grilltemp').text(data['grilltemp']);
  $('.grilltarget').text(data['grilltarget']);
  $('.probetemp').text(data['probetemp']);
  $('.probetarget').text(data['probetarget']);
  $('.grilltemp').toggleClass('yft');
  $('.probetemp').toggleClass('yft');
  $('#datetime').text("Last Checked at: " + (new Date()))
  setTimeout(function () {
    $('.grilltemp').toggleClass('yft');
    $('.probetemp').toggleClass('yft')
  }, 2000)
})
.fail(function() {
  console.log("error");
})
.always(function() {
  console.log("complete");
});
}
function btoc() {
  $.ajax({
    url: '/cmd',
    type: 'POST',
    dataType: 'json',
    data: '{"cmd": "btoc"}'
  })
  .done(function() {
    console.log("success");
  })
  .fail(function() {
    console.log("error");
  })
  .always(function() {
    console.log("complete");
  });
}
function poweron() {
  $.ajax({
    url: '/power',
    type: 'POST',
    dataType: 'json',
    data: '{"cmd": "on"}'
  })
  .done(function() {
    console.log("success");
  })
  .fail(function() {
    console.log("error");
  })
  .always(function() {
    console.log("complete");
  });
}
function poweroff() {
  $.ajax({
    url: '/power',
    type: 'POST',
    dataType: 'json',
    data: '{"cmd": "off"}'
  })
  .done(function() {
    console.log("success");
  })
  .fail(function() {
    console.log("error");
  })
  .always(function() {
    console.log("complete");
  });
}
function setgrilltarget() {
  $.ajax({
    url: 'temp/grilltarget',
    type: 'POST',
    dataType: 'json',
    data: '{"grill": '+$("#grillTargetTemp").val()+'}'
  })
  .done(function() {
    console.log("success");
  })
  .fail(function() {
    console.log("error");
  })
  .always(function() {
    console.log("complete");
  });

}
function setprobetarget() {
  $.ajax({
    url: 'temp/probetarget',
    type: 'POST',
    dataType: 'json',
    data: '{"probe": '+$("#probeTargetTemp").val()+'}'
  })
  .done(function() {
    console.log("success");
  })
  .fail(function() {
    console.log("error");
  })
  .always(function() {
    console.log("complete");
  });

}

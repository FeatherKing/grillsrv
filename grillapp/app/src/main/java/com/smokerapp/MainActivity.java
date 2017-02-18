package com.smokerapp;

import android.content.Intent;
import android.net.Uri;
import android.os.AsyncTask;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.NumberPicker;
import android.widget.TextView;
import com.google.android.gms.appindexing.Action;
import com.google.android.gms.appindexing.AppIndex;
import com.google.android.gms.common.api.GoogleApiClient;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.List;

public class MainActivity extends AppCompatActivity {

    ImageButton ib;
    Button button;
    TextView results;
    Button setButton;
    NumberPicker np;
    /**
     * ATTENTION: This was auto-generated to implement the App Indexing API.
     * See https://g.co/AppIndexing/AndroidStudio for more information.
     */
    private GoogleApiClient client;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);


//        Button settingsButton = (Button) findViewById(R.id.settings_button);
//        settingsButton.setOnClickListener(new View.OnClickListener() {
//            @Override
//            public void onClick(View view) {
//                Intent intent = new Intent(MainActivity.this, Main2Activity.class);
//                startActivity(intent);
//            }
//        });
        ib = (ImageButton) findViewById(R.id.settings_button);
        ib.setOnClickListener(ibLis);


        button = (Button) findViewById(R.id.queryButton);
        button.setOnClickListener(ibLis2);
        results = (TextView) findViewById(R.id.responseView);

        setButton = (Button) findViewById(R.id.setgrilltempButton);
        setButton.setOnClickListener(setGrillTempLis);

        //spinner
        String[] nums = new String[501];

        for(int i=0; i<nums.length; i++)
            nums[i] = Integer.toString(i*1);
        np = (NumberPicker) findViewById(R.id.np);
        np.setMaxValue(nums.length-1);
        np.setMinValue(0);
        np.setWrapSelectorWheel(false);
        np.setDisplayedValues(nums);





        // ATTENTION: This was auto-generated to implement the App Indexing API.
        // See https://g.co/AppIndexing/AndroidStudio for more information.
        client = new GoogleApiClient.Builder(this).addApi(AppIndex.API).build();
    }

    private View.OnClickListener ibLis = new View.OnClickListener() {

        @Override
        public void onClick(View v) {
            // TODO Auto-generated method stub
            //START YOUR ACTIVITY HERE AS
            Intent intent = new Intent(MainActivity.this, SettingsActivity.class);
            startActivity(intent);
        }
    };

    private View.OnClickListener ibLis2 = new View.OnClickListener() {

        @Override
        public void onClick(View v)  {
            // TODO Auto-generated method stub
            //START YOUR ACTIVITY HERE AS
            results.setText("YAYAYAYAYAY");

            new SendGetRequest().execute();



        }
    };

    private View.OnClickListener setGrillTempLis = new View.OnClickListener() {

        @Override
        public void onClick(View v) {
            // TODO Auto-generated method stub
            //START YOUR ACTIVITY HERE AS

            new SendPostRequest().execute(np.getValue());
        }
    };

    @Override
    public void onStart() {
        super.onStart();

        // ATTENTION: This was auto-generated to implement the App Indexing API.
        // See https://g.co/AppIndexing/AndroidStudio for more information.
        client.connect();
        Action viewAction = Action.newAction(
                Action.TYPE_VIEW, // TODO: choose an action type.
                "Main Page", // TODO: Define a title for the content shown.
                // TODO: If you have web page content that matches this app activity's content,
                // make sure this auto-generated web page URL is correct.
                // Otherwise, set the URL to null.
                Uri.parse("http://host/path"),
                // TODO: Make sure this auto-generated app URL is correct.
                Uri.parse("android-app://com.smokerapp/http/host/path")
        );
        AppIndex.AppIndexApi.start(client, viewAction);
    }

    @Override
    public void onStop() {
        super.onStop();

        // ATTENTION: This was auto-generated to implement the App Indexing API.
        // See https://g.co/AppIndexing/AndroidStudio for more information.
        Action viewAction = Action.newAction(
                Action.TYPE_VIEW, // TODO: choose an action type.
                "Main Page", // TODO: Define a title for the content shown.
                // TODO: If you have web page content that matches this app activity's content,
                // make sure this auto-generated web page URL is correct.
                // Otherwise, set the URL to null.
                Uri.parse("http://host/path"),
                // TODO: Make sure this auto-generated app URL is correct.
                Uri.parse("android-app://com.smokerapp/http/host/path")
        );
        AppIndex.AppIndexApi.end(client, viewAction);
        client.disconnect();
    }

    class SendGetRequest extends AsyncTask<Void, Void, String> {

        protected void onPreExecute() {
//            progressBar.setVisibility(View.VISIBLE);
//            responseView.setText("");
        }

        protected String doInBackground(Void... urls) {
//            String email = emailText.getText().toString();
            // Do some validation here

            try {
                URL url = new URL("http://24.11.74.43:9999/info");
                HttpURLConnection urlConnection = (HttpURLConnection) url.openConnection();
                try {
                    BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(urlConnection.getInputStream()));
                    StringBuilder stringBuilder = new StringBuilder();
                    String line;
                    while ((line = bufferedReader.readLine()) != null) {
                        stringBuilder.append(line).append("\n");
                    }
                    bufferedReader.close();
                    return stringBuilder.toString();
                } finally {
                    urlConnection.disconnect();
                }
            } catch (Exception e) {
                Log.e("ERROR", e.getMessage(), e);
                return null;
            }
        }

        protected void onPostExecute(String response) {
            if (response == null) {
                response = "THERE WAS AN ERROR";
            }
//            progressBar.setVisibility(View.GONE);
            Log.i("INFO", response);
            results.setText(response);
        }
    }

    private class SendPostRequest extends AsyncTask<Integer,Integer,String> {

//        public SendPostRequest(int temptarget) {
//            super();
//            // do stuff
//        }

        protected void onPreExecute() {
//            progressBar.setVisibility(View.VISIBLE);
//            responseView.setText("");
        }

        protected String doInBackground(Integer... params) {
            URL url;
            HttpURLConnection conn = null;
            JSONObject jsonObject = new JSONObject();
            OutputStream os;

            try {
                url = new URL("http://24.11.74.43:9999/temp/grilltarget");
                conn = (HttpURLConnection) url.openConnection();
                conn.setRequestProperty("Content-Type", "application/json");
                conn.setDoInput(true);
                conn.setDoOutput(true);
                conn.setRequestMethod("POST");

                try {
                    jsonObject.put("grill", Integer.valueOf(params[0]));
                } catch (JSONException e) {
                    e.printStackTrace();
                }

                os = conn.getOutputStream();
                os.write(jsonObject.toString().getBytes());
                os.flush();
                conn.getResponseCode();
                return String.valueOf(conn.getResponseCode());
            } catch (Exception e) {
                return new String("Exception: " + e.getMessage());
            } finally {
                conn.disconnect();
            }

        }



        protected void onPostExecute(String response) {
            if (response == null) {
                response = "THERE WAS AN ERROR";
            }
//            progressBar.setVisibility(View.GONE);
            Log.i("INFO", response);
            results.setText(response);  //Show the response in this area on the screen
        }
    }

}














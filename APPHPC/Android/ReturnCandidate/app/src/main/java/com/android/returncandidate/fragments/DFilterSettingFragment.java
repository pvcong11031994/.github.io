package com.android.returncandidate.fragments;

import android.app.Dialog;
import android.content.DialogInterface;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.support.v4.app.FragmentManager;
import android.view.KeyEvent;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.Button;
import android.widget.CompoundButton;
import android.widget.ImageButton;
import android.widget.TextView;
import android.widget.ToggleButton;

import com.android.returncandidate.R;
import com.android.returncandidate.common.constants.Constants;
import com.android.returncandidate.common.utils.Common;
import com.android.returncandidate.common.utils.FlagSettingNew;
import com.android.returncandidate.common.utils.FlagSettingOld;

import com.android.returncandidate.common.utils.ProcessDialogSetting;
import com.android.returncandidate.common.utils.ProcessDialogUpdate;

import java.util.ArrayList;

/**
 * <h1>List select Dialog</h1>
 * <p>
 * Dialog list select screen
 *
 * @author cong-pv
 * @since 2018/07/06.
 */
@SuppressWarnings("deprecation")
public class DFilterSettingFragment extends DialogFragment implements View.OnClickListener {

    /**
     * Interface to item selected to activity
     */
    public interface ItemSelectedFilterSettingDialogListener {

        /**
         * Function send list selected data to activity
         */
        void onLitSelectedFilterSettingDialog(FlagSettingNew flagSettingNew, FlagSettingOld flagSettingOld,
                                              Boolean flagFilterSubmit, String flagSwitchOCR);
    }

    /**
     * Button back
     */
    private ImageButton imbBack;

    private FlagSettingNew flagSettingNew;
    private FlagSettingOld flagSettingOld;
    private String flagSwitchOCR;
    private Common common;

    //Check submit filter
    private Boolean flagFilterSubmit;

    //Check yes/no joubi
    private String flagJoubi = Constants.VALUE_YES_STANDING;
    Button btnClassification, btnPublisher, btnReleaseDate, btnUndisturbed,
            btnNumberOfStocks, btnStocksPercent, btnDecision, btnNotCovered;
    TextView txvClassification, txvPublisher, txvReleaseDate, txvUndisturbed,
            txvNumberOfStocks, txvStocksPercent;
    ToggleButton tbNotCovered;


    /**
     * Init Dialog Product Detail with custom layout
     *
     * @param container         is view group
     * @param inflater          is layout inflater
     * @param saveInstanceState bundle of dialog
     * @return show dialog
     */
    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle saveInstanceState) {


        //Default filter submit false
        flagFilterSubmit = false;

        // Init custom product detail layout
        View rootView = inflater.inflate(R.layout.fragement_setting, container, false);

        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }
        flagSettingNew = new FlagSettingNew();
        flagSettingOld = new FlagSettingOld();
        common = new Common();
        if (getArguments() != null) {
            common.SetArgumentsFlagData(flagSettingNew, flagSettingOld, getArguments());
            flagSwitchOCR = getArguments().getString(Constants.FLAG_SWITCH_OCR);
        }

        // Variable button click
        btnClassification = (Button) rootView.findViewById(R.id.btn_classification);
        btnPublisher = (Button) rootView.findViewById(R.id.btn_publisher);
        btnReleaseDate = (Button) rootView.findViewById(R.id.btn_release_date);
        btnUndisturbed = (Button) rootView.findViewById(R.id.btn_undisturbed);
        btnNumberOfStocks = (Button) rootView.findViewById(R.id.btn_number_of_stocks);
        btnStocksPercent = (Button) rootView.findViewById(R.id.btn_stocks_percent);
        btnNotCovered = (Button) rootView.findViewById(R.id.btn_not_covered);
        imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);
        btnDecision = (Button) rootView.findViewById(R.id.btn_decision);

        //Variable text view
        txvClassification = (TextView) rootView.findViewById(R.id.txv_classification);
        txvPublisher = (TextView) rootView.findViewById(R.id.txv_publisher);
        txvReleaseDate = (TextView) rootView.findViewById(R.id.txv_release_date);
        txvUndisturbed = (TextView) rootView.findViewById(R.id.txv_undisturbed);
        txvNumberOfStocks = (TextView) rootView.findViewById(R.id.txv_number_of_stocks);
        txvStocksPercent = (TextView) rootView.findViewById(R.id.txv_stocks_percent);

        //ToggleButton
        tbNotCovered = (ToggleButton) rootView.findViewById(R.id.tb_not_covered);

        //Check if group2 len = 2 and select all
        int poisitionSelectAll = -1;
        if (flagSettingNew.getFlagClassificationGroup2Cd().size() == 2) {
            for (int i = 0; i < flagSettingNew.getFlagClassificationGroup2Cd().size(); i++) {
                if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup2Cd().get(i))) {
                    poisitionSelectAll = i;
                }
            }
        }

        //Set text view
        //Check show classification
        if (flagSettingNew.getFlagClassificationGroup1Cd().size() == 1) {
            if (flagSettingNew.getFlagClassificationGroup2Cd().size() == 1) {
                if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup2Cd().get(0))) {
                    txvClassification.setText(Constants.VALUE_STR_WHEN_SELECT_ALL);
                } else {
                    txvClassification.setText(flagSettingNew.getFlagClassificationGroup1Name().get(0) + " (" + flagSettingNew.getFlagClassificationGroup2Name().get(0) + ")");
                }
            } else if (poisitionSelectAll != -1) {
                txvClassification.setText(flagSettingNew.getFlagClassificationGroup1Name().get(0) + " (" + flagSettingNew.getFlagClassificationGroup2Name().get(1 - poisitionSelectAll) + ")");
            } else {
                txvClassification.setText(flagSettingNew.getFlagClassificationGroup1Name().get(0) + " (" + Constants.VALUE_STR_WHEN_SELECT_MULTI + ")");
            }
        } else if (flagSettingNew.getFlagClassificationGroup1Cd().size() >= 1) {
            if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagClassificationGroup1Cd().get(0))) {
                txvClassification.setText(Constants.VALUE_STR_WHEN_SELECT_ALL);
            } else {
                txvClassification.setText(Constants.VALUE_STR_WHEN_SELECT_MULTI);
            }
        } else {
            txvClassification.setText(Constants.VALUE_STR_WHEN_SELECT_ALL);
        }


        //Check show publisher
        if (flagSettingNew.getFlagPublisherShowScreen().size() == Constants.VALUE_WHEN_SELECT_ONE) {
            txvPublisher.setText(flagSettingNew.getFlagPublisherShowScreen().get(0));
        } else if (flagSettingNew.getFlagPublisherShowScreen().size() > Constants.VALUE_WHEN_SELECT_ONE) {
            if (Constants.ID_ROW_ALL.equals(flagSettingNew.getFlagPublisher().get(0))) {
                txvPublisher.setText(Constants.ROW_ALL);
            } else {
                txvPublisher.setText(Constants.VALUE_STR_WHEN_SELECT_MULTI);
            }
        } else {
            txvPublisher.setText(Constants.ROW_ALL);
        }
        //Show Other
        txvReleaseDate.setText(flagSettingNew.getFlagReleaseDateShowScreen());
        txvUndisturbed.setText(flagSettingNew.getFlagUndisturbedShowScreen());
        txvNumberOfStocks.setText(flagSettingNew.getFlagNumberOfStocksShowScreen());
        txvStocksPercent.setText(flagSettingNew.getFlagStockPercentShowScreen());

        //Event click button & text view setting
        btnClassification.setOnClickListener(this);
        txvClassification.setOnClickListener(this);
        btnPublisher.setOnClickListener(this);
        txvPublisher.setOnClickListener(this);
        btnReleaseDate.setOnClickListener(this);
        txvReleaseDate.setOnClickListener(this);
        btnUndisturbed.setOnClickListener(this);
        txvUndisturbed.setOnClickListener(this);
        btnNumberOfStocks.setOnClickListener(this);
        txvNumberOfStocks.setOnClickListener(this);
        btnStocksPercent.setOnClickListener(this);
        txvStocksPercent.setOnClickListener(this);
        btnDecision.setOnClickListener(this);
        imbBack.setOnClickListener(this);
        btnNotCovered.setOnClickListener(this);
        if (Constants.VALUE_NO_STANDING.equals(flagSettingNew.getFlagJoubi())) {
            tbNotCovered.setChecked(false);
            tbNotCovered.setText(Constants.NO_STANDING);
        } else {
            tbNotCovered.setChecked(true);
            tbNotCovered.setText(Constants.YES_STANDING);
        }

        tbNotCovered.setOnCheckedChangeListener(new CompoundButton.OnCheckedChangeListener() {
            @Override
            public void onCheckedChanged(CompoundButton buttonView, boolean isChecked) {
                showToggleButton();
            }
        });
        //Event lick back device
        getDialog().setOnKeyListener(new DialogInterface.OnKeyListener() {
            @Override
            public boolean onKey(DialogInterface dialog, int keyCode, KeyEvent event) {
                if (keyCode == KeyEvent.KEYCODE_BACK) {
                    backDialogReturn();
                    return true;
                }
                return false;
            }
        });
        return rootView;
    }

    @Override
    public void onPrepareOptionsMenu(Menu menu) {
        super.onPrepareOptionsMenu(menu);
    }

    @Override
    public void onResume() {
        super.onResume();

        ArrayList<FlagSettingNew> arrFlagSettingNew = new ArrayList<>();
        arrFlagSettingNew.add(flagSettingNew);

        new ProcessDialogSetting(getActivity(), getView()).execute(arrFlagSettingNew);
    }

    /**
     * Set size layout on start
     */
    @Override
    public void onStart() {
        super.onStart();
        Dialog dialog = getDialog();
        if (dialog != null) {
            if (dialog.getWindow() != null) {
                dialog.getWindow().setLayout(ViewGroup.LayoutParams.FILL_PARENT,
                        ViewGroup.LayoutParams.FILL_PARENT);
            }
        }
    }

    /**
     * On click
     *
     * @param v is View on click listener
     */
    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            //Event click button 分類
            case R.id.btn_classification:
                showScreenGroup1Cd();
                break;
            //Event click text view 分類
            case R.id.txv_classification:
                showScreenGroup1Cd();
                break;
            //Event click button 出版社
            case R.id.btn_publisher:
                showScreenPublisher();
                break;
            //Event click text view 出版社
            case R.id.txv_publisher:
                showScreenPublisher();
                break;
            //Event click button 発売日
            case R.id.btn_release_date:
                showScreenReleaseDate();
                break;
            //Event click text view 発売日
            case R.id.txv_release_date:
                showScreenReleaseDate();
                break;
            //Event click button 未動期間
            case R.id.btn_undisturbed:
                showScreenUndisturbed();
                break;
            //Event click text view 未動期間
            case R.id.txv_undisturbed:
                showScreenUndisturbed();
                break;
            //Event click button 在庫数
            case R.id.btn_number_of_stocks:
                showScreenNumberOfStocks();
                break;
            //Event click text view 在庫数
            case R.id.txv_number_of_stocks:
                showScreenNumberOfStocks();
                break;
            //Event click button 在庫％
            case R.id.btn_stocks_percent:
                showScreenStocksPercent();
                break;
            //Event click text view 在庫％
            case R.id.txv_stocks_percent:
                showScreenStocksPercent();
                break;
            //Event click button 対象外
            //case R.id.btn_not_covered:
            //    break;
            case R.id.imb_back:
                flagFilterSubmit = false;
                backDialogReturn();
                break;
            case R.id.btn_decision:
                flagFilterSubmit = true;
                submitFilter();
                break;
        }
    }

    /**
     * Function submit filter
     */

    private void submitFilter() {

        flagSettingNew.setFlagJoubi(flagSettingNew.getFlagJoubi());

        // Save flag new
        saveFlagNewToFlagOld();

        ArrayList<FlagSettingNew> arrFlagSettingNew = new ArrayList<>();
        arrFlagSettingNew.add(flagSettingNew);

        new ProcessDialogUpdate(getActivity(), flagSwitchOCR).execute(arrFlagSettingNew);

        // move to select year dialog
        ItemSelectedFilterSettingDialogListener activity = (ItemSelectedFilterSettingDialogListener) getActivity();
        activity.onLitSelectedFilterSettingDialog(flagSettingNew, flagSettingOld, flagFilterSubmit, flagSwitchOCR);
        dismiss();
    }

    //Save flag new into flag old
    private void saveFlagNewToFlagOld() {

        flagSettingOld.setFlagClassificationGroup1Cd(flagSettingNew.getFlagClassificationGroup1Cd());
        flagSettingOld.setFlagClassificationGroup1Name(flagSettingNew.getFlagClassificationGroup1Name());
        flagSettingOld.setFlagClassificationGroup2Cd(flagSettingNew.getFlagClassificationGroup2Cd());
        flagSettingOld.setFlagClassificationGroup2Name(flagSettingNew.getFlagClassificationGroup2Name());
        //save flag publisher
        flagSettingOld.setFlagPublisher(flagSettingNew.getFlagPublisher());
        flagSettingOld.setFlagPublisherShowScreen(flagSettingNew.getFlagPublisherShowScreen());
        //save flag release date
        flagSettingOld.setFlagReleaseDate(flagSettingNew.getFlagReleaseDate());
        flagSettingOld.setFlagReleaseDateShowScreen(flagSettingNew.getFlagReleaseDateShowScreen());
        //save flag undisturbed
        flagSettingOld.setFlagUndisturbed(flagSettingNew.getFlagUndisturbed());
        flagSettingOld.setFlagUndisturbedShowScreen(flagSettingNew.getFlagUndisturbedShowScreen());
        //save flag number of stocks
        flagSettingOld.setFlagNumberOfStocks(flagSettingNew.getFlagNumberOfStocks());
        flagSettingOld.setFlagNumberOfStocksShowScreen(flagSettingNew.getFlagNumberOfStocksShowScreen());
        //save flag stocks percent
        flagSettingOld.setFlagStockPercent(flagSettingNew.getFlagStockPercent());
        flagSettingOld.setFlagStockPercentShowScreen(flagSettingNew.getFlagStockPercentShowScreen());
        //save flag joubi
        flagSettingOld.setFlagJoubi(flagSettingNew.getFlagJoubi());

    }

    /**
     * Function call when click button 分類 (classification) -> move screen select group1_cd
     */

    private void showScreenGroup1Cd() {

        //Call fragment group 1 cd
        DFilterListGenreFragment dFilterListGenreFragment = new DFilterListGenreFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        //Set bundle to argument
        dFilterListGenreFragment.setArguments(bundle);
        dFilterListGenreFragment.show(fm, null);
        dismiss();
    }

    /**
     * Function call when click toggle button
     */

    private void showToggleButton() {

        if (!tbNotCovered.isChecked()) {
            flagJoubi = Constants.VALUE_NO_STANDING;
        } else {
            flagJoubi = Constants.VALUE_YES_STANDING;
        }
        flagSettingNew.setFlagJoubi(flagJoubi);
        //Restart fragment
        onResume();
    }

    /**
     * Function call when click button 出版社 (publisher) -> move screen select publisher
     */

    private void showScreenPublisher() {

        //Call fragment publisher
        DFilterListPublisherFragment dFilterListPublisherFragment = new DFilterListPublisherFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        //Set bundle to argument
        dFilterListPublisherFragment.setArguments(bundle);
        dFilterListPublisherFragment.show(fm, null);
        dismiss();
    }

    /**
     * Function call when click button 発売日 (release date) -> move screen select release date
     */

    private void showScreenReleaseDate() {

        //Call fragment
        DFilterListMonthReleaseDateFragment dFilterListMonthReleaseDateFragment = new DFilterListMonthReleaseDateFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        //Set bundle to argument
        dFilterListMonthReleaseDateFragment.setArguments(bundle);
        dFilterListMonthReleaseDateFragment.show(fm, null);
        dismiss();
    }

    /**
     * Function call when click button 未動期間 (undisturbed) -> move screen select undisturbed
     */

    private void showScreenUndisturbed() {

        //Call fragment
        DFilterListMonthUndisturbedFragment dFilterListMonthUndisturbedFragment = new DFilterListMonthUndisturbedFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        //Set bundle to argument
        dFilterListMonthUndisturbedFragment.setArguments(bundle);
        dFilterListMonthUndisturbedFragment.show(fm, null);
        dismiss();
    }

    /**
     * Function call when click button 在庫数 (number of stocks) -> move screen select number of stocks
     */
    private void showScreenNumberOfStocks() {

        //Call fragment
        DFilterListNumberOfStockFragment dFilterListNumberOfStockFragment = new DFilterListNumberOfStockFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        //Set bundle to argument
        dFilterListNumberOfStockFragment.setArguments(bundle);
        dFilterListNumberOfStockFragment.show(fm, null);
        dismiss();
    }

    /**
     * Function call when click button 在庫％ (stocks percent) -> move screen select stocks percent
     */
    private void showScreenStocksPercent() {

        //Call fragment
        DFilterListStocksPercentFragment dFilterListStocksPercentFragment = new DFilterListStocksPercentFragment();
        FragmentManager fm = getActivity().getSupportFragmentManager();
        Bundle bundle = common.DataPutActivity(flagSettingNew, flagSettingOld);
        //put flag switch OCR
        bundle.putString(Constants.FLAG_SWITCH_OCR, flagSwitchOCR);

        //Set bundle to argument
        dFilterListStocksPercentFragment.setArguments(bundle);
        dFilterListStocksPercentFragment.show(fm, null);
        dismiss();
    }

    /**
     * Back dialog setting
     */
    private void backDialogReturn() {

        //back flag setting joubi
        flagSettingNew.setFlagJoubi(flagSettingOld.getFlagJoubi());

        // move to select year dialog
        ItemSelectedFilterSettingDialogListener activity = (ItemSelectedFilterSettingDialogListener) getActivity();
        activity.onLitSelectedFilterSettingDialog(flagSettingNew, flagSettingOld, flagFilterSubmit, flagSwitchOCR);
        dismiss();
    }

}
